//
//  ContentView.swift
//  btfs-test
//
//  Created by Talha on 11/04/22.
//

import SwiftUI

struct ContentView: View {
    @State private var txt: String = "daemon --chain-id 199"

        var body: some View {
            VStack{
                TextField("daemon --chain-id 199", text: $txt)
                .textFieldStyle(RoundedBorderTextFieldStyle())
                Button("DO cli call"){
                    let str = mainC(UnsafeMutablePointer<Int8>(mutating: (self.txt as NSString).utf8String))
                    self.txt = String.init(cString: str!, encoding: .utf8)!
                    // don't forget to release the memory to the C String
                    str?.deallocate()
                }

                Spacer()
            }
            .padding(.all, 15)
        }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        Group {
            ContentView()
                .previewInterfaceOrientation(.portrait)
            ContentView()
                .previewInterfaceOrientation(.portrait)
        }
    }
}

