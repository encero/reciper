//
//  receptikyApp.swift
//  receptiky
//
//  Created by Matous Michalik on 05.03.2022.
//

import SwiftUI

@main
struct receptikyApp: App {
    
    @State var configured = Settings.shared.graphqlServerURL != nil
    
    
    
    var body: some Scene {
        WindowGroup {
            if configured {
                ContentView()
                    .environmentObject(RecipeDataManager())
            } else {
                SettingsView() {
                    configured = true
                }
            }
        }
    }
}
