package commands

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	cmdenv "github.com/bittorrent/go-btfs/core/commands/cmdenv"

	cmds "github.com/bittorrent/go-btfs-cmds"
	unixfs "github.com/bittorrent/go-unixfs"
	unixfs_pb "github.com/bittorrent/go-unixfs/pb"
	iface "github.com/bittorrent/interface-go-btfs-core"
	options "github.com/bittorrent/interface-go-btfs-core/options"
	path "github.com/bittorrent/interface-go-btfs-core/path"
)

// LsLink contains printable data for a single ipld link in ls output
type LsLink struct {
	Name, Hash string
	Size       uint64
	Type       unixfs_pb.Data_DataType
	Target     string
	Mode       os.FileMode
	Mtime      time.Time
}

// LsObject is an element of LsOutput
// It can represent all or part of a directory
type LsObject struct {
	Hash  string
	Links []LsLink
}

// LsOutput is a set of printable data for directories,
// it can be complete or partial
type LsOutput struct {
	Objects []LsObject
}

const (
	lsHeadersOptionNameTime = "headers"
	lsResolveTypeOptionName = "resolve-type"
	lsSizeOptionName        = "size"
	lsStreamOptionName      = "stream"
	lsMTimeOptionName       = "mtime"
	lsModeOptionName        = "mode"
)

var LsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List directory contents for Unix filesystem objects.",
		ShortDescription: `
Displays the contents of an BTFS or BTNS object(s) at the given path, with
the following format:

  <link base58 hash> <link size in bytes> <link name>

The JSON output contains type information.
`,
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("btfs-path", true, true, "The path to the BTFS object(s) to list links from.").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.BoolOption(lsHeadersOptionNameTime, "v", "Print table headers (Hash, Size, Name)."),
		cmds.BoolOption(lsResolveTypeOptionName, "Resolve linked objects to find out their types.").WithDefault(true),
		cmds.BoolOption(lsSizeOptionName, "Resolve linked objects to find out their file size.").WithDefault(true),
		cmds.BoolOption(lsStreamOptionName, "s", "Enable experimental streaming of directory entries as they are traversed."),
		cmds.BoolOption(lsMTimeOptionName, "t", "Print modification time."),
		cmds.BoolOption(lsModeOptionName, "m", "Print mode."),
	},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		api, err := cmdenv.GetApi(env, req)
		if err != nil {
			return err
		}

		resolveType, _ := req.Options[lsResolveTypeOptionName].(bool)
		resolveSize, _ := req.Options[lsSizeOptionName].(bool)
		stream, _ := req.Options[lsStreamOptionName].(bool)

		err = req.ParseBodyArgs()
		if err != nil {
			return err
		}
		paths := req.Arguments

		enc, err := cmdenv.GetCidEncoder(req)
		if err != nil {
			return err
		}

		var processLink func(path string, link LsLink) error
		var dirDone func(i int)

		processDir := func() (func(path string, link LsLink) error, func(i int)) {
			return func(path string, link LsLink) error {
				output := []LsObject{{
					Hash:  path,
					Links: []LsLink{link},
				}}
				return res.Emit(&LsOutput{output})
			}, func(i int) {}
		}
		done := func() error { return nil }

		if !stream {
			output := make([]LsObject, len(req.Arguments))

			processDir = func() (func(path string, link LsLink) error, func(i int)) {
				// for each dir
				outputLinks := make([]LsLink, 0)
				return func(path string, link LsLink) error {
						// for each link
						outputLinks = append(outputLinks, link)
						return nil
					}, func(i int) {
						// after each dir
						sort.Slice(outputLinks, func(i, j int) bool {
							return outputLinks[i].Name < outputLinks[j].Name
						})

						output[i] = LsObject{
							Hash:  paths[i],
							Links: outputLinks,
						}
					}
			}

			done = func() error {
				return cmds.EmitOnce(res, &LsOutput{output})
			}
		}

		for i, fpath := range paths {
			results, err := api.Unixfs().Ls(req.Context, path.New(fpath),
				options.Unixfs.ResolveChildren(resolveSize || resolveType))
			if err != nil {
				return err
			}

			processLink, dirDone = processDir()
			for link := range results {
				if link.Err != nil {
					return link.Err
				}
				var ftype unixfs_pb.Data_DataType
				switch link.Type {
				case iface.TFile:
					ftype = unixfs.TFile
				case iface.TDirectory:
					ftype = unixfs.TDirectory
				case iface.TSymlink:
					ftype = unixfs.TSymlink
				}
				lsLink := LsLink{
					Name: link.Name,
					Hash: enc.Encode(link.Cid),

					Size:   link.Size,
					Type:   ftype,
					Target: link.Target,
					Mode:   link.Mode,
					Mtime:  link.ModTime,
				}
				if err := processLink(paths[i], lsLink); err != nil {
					return err
				}
			}
			dirDone(i)
		}
		return done()
	},
	PostRun: cmds.PostRunMap{
		cmds.CLI: func(res cmds.Response, re cmds.ResponseEmitter) error {
			req := res.Request()
			lastObjectHash := ""

			for {
				v, err := res.Next()
				if err != nil {
					if err == io.EOF {
						return nil
					}
					return err
				}
				out := v.(*LsOutput)
				lastObjectHash = tabularOutput(req, os.Stdout, out, lastObjectHash, false)
			}
		},
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, out *LsOutput) error {
			// when streaming over HTTP using a text encoder, we cannot render breaks
			// between directories because we don't know the hash of the last
			// directory encoder
			ignoreBreaks, _ := req.Options[lsStreamOptionName].(bool)
			tabularOutput(req, w, out, "", ignoreBreaks)
			return nil
		}),
	},
	Type: LsOutput{},
}

func tabularOutput(req *cmds.Request, w io.Writer, out *LsOutput, lastObjectHash string, ignoreBreaks bool) string {
	headers, _ := req.Options[lsHeadersOptionNameTime].(bool)
	stream, _ := req.Options[lsStreamOptionName].(bool)
	size, _ := req.Options[lsSizeOptionName].(bool)
	mtime, _ := req.Options[lsMTimeOptionName].(bool)
	mode, _ := req.Options[lsModeOptionName].(bool)
	// in streaming mode we can't automatically align the tabs
	// so we take a best guess
	var minTabWidth int
	if stream {
		minTabWidth = 10
	} else {
		minTabWidth = 1
	}

	multipleFolders := len(req.Arguments) > 1

	tw := tabwriter.NewWriter(w, minTabWidth, 2, 1, ' ', 0)

	for _, object := range out.Objects {

		if !ignoreBreaks && object.Hash != lastObjectHash {
			if multipleFolders {
				if lastObjectHash != "" {
					fmt.Fprintln(tw)
				}
				fmt.Fprintf(tw, "%s:\n", object.Hash)
			}
			if headers {
				s := "Hash\tName"
				if size {
					s = "Hash\tSize\tName"
				}

				s = buildHeader(mode, "Mode", s)
				s = buildHeader(mtime, "Mtime", s)

				fmt.Fprintln(tw, s)
			}
			lastObjectHash = object.Hash
		}

		for _, link := range object.Links {
			var s string
			switch link.Type {
			case unixfs.TDirectory, unixfs.THAMTShard, unixfs.TMetadata:
				if size {
					s = "%[1]s\t-\t%[3]s/"
				} else {
					s = "%[1]s\t%[3]s/"
				}
				s = buildString(mode, s, 4)
				s = buildString(mtime, s, 5)
				s = s + "\n"
			default:
				if size {
					s = "%[1]s\t%[2]v\t%[3]s"
				} else {
					s = "%[1]s\t%[3]s"
				}
				s = buildString(mode, s, 4)
				s = buildString(mtime, s, 5)
				s = s + "\n"
			}

			modeS := "-"
			mtimeS := "-"

			if link.Mode != 0 {
				modeS = link.Mode.String()
			}
			if link.Mtime.Unix() != 0 {
				mtimeS = link.Mtime.Format("2 Jan 2006, 15:04:05 MST")
			}
			fmt.Fprintf(tw, s, link.Hash, link.Size, link.Name, modeS, mtimeS)
		}
	}
	tw.Flush()
	return lastObjectHash
}

func buildString(set bool, s string, index int) string {
	if set {
		return fmt.Sprintf("%s\t%%[%d]s", s, index)
	}
	return s
}

func buildHeader(set bool, name, s string) string {
	if set {
		return s + "\t" + name
	}
	return s
}
