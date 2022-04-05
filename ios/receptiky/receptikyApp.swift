//
//  receptikyApp.swift
//  receptiky
//
//  Created by Matous Michalik on 05.03.2022.
//

import SwiftUI

@main
struct receptikyApp: App {
    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(RecipeDataManager())
        }
    }
}
