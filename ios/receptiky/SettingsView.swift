//
//  InitialSetupView.swift
//  receptiky
//
//  Created by Matous Michalik on 13.04.2022.
//

import SwiftUI

struct SettingsView: View {
    @State var url: String = Settings.shared.graphqlServerURL ?? ""

    @State var password = Settings.shared.graphqlServerPassword ?? ""
    @State var username = Settings.shared.graphqlServerUsername ?? ""
    
    @Environment(\.dismiss) private var dismiss
    
    @State var error: String?
    
    @State var useAuth = Settings.shared.graphqlServerPassword != nil
    
    var callback: (()->())?
    
    var body:
    some View {
        if callback != nil {
            Text("Initial application setup").font(.title).padding(.top, 10)
        }
        if error != nil {
            Text(error!)
                .tint(.red)
                .padding(10)
        }
        Form {
            Section("Configuration") {
                TextField("Server url", text:$url)
                    .autocapitalization(.none)
                    .textContentType(.URL)
                Toggle("Use Auth?", isOn: $useAuth.animation())
                if useAuth {
                    TextField("Username", text:$username )
                        .autocapitalization(.none)
                    SecureField("Password", text: $password)
                        .autocapitalization(.none)
                }
            }
            
            
            Button("Save") {
                if valid {
                    error = nil
                    
                    let reqURL = URL(string: url)
                    if reqURL == nil {
                        error = "Invalid url"
                        return
                    }
                    
                    var request = URLRequest(url: reqURL!)
                    
                    request.httpMethod = "GET"
                    
                    if useAuth {
                        let headerValue = Data("\(username):\(password)".utf8).base64EncodedString()
                        request.setValue("Basic \(headerValue)", forHTTPHeaderField: "Authentication")
                    }
                    
                    
                    let task = URLSession.shared.dataTask(with: request) { data, response, error in
                        print("url called")
                        if error != nil {
                            self.error = "Cant reach server: \(error!.localizedDescription)"
                            return
                        }
                        
                        // todo: create status gql query and check it here
                        
                        Settings.shared.graphqlServerURL = url
                        Settings.shared.graphqlServerUsername = username
                        Settings.shared.graphqlServerPassword = password
                        
                        if callback != nil {
                            callback!()
                        } else {
                            dismiss()
                        }
                    }
                    
                    task.resume()
                } else {
                    print("not valid")
                }
            }.disabled(!valid)
        }.navigationTitle("Settings")
    }
    
    var valid: Bool {
        if !url.starts(with: "https://") {
            return false
        }
        
        if useAuth && (password.isEmpty || username.isEmpty) {
            return false
        }
        
        return true
    }
    
}
