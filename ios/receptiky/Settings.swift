//
//  Settings.swift
//  receptiky
//
//  Created by Matous Michalik on 13.04.2022.
//

import Foundation

class Settings {
    private static let keyGraphQLServerURL = "gqlURL"
    private static let keyGraphQLServerPassword = "gqlPassword"
    private static let keyGraphQLServerUsername = "gqlUsername"
    
    static var shared = Settings()
    
    private init() {
        
    }
    
    var graphqlServerURL: String? {
        get {
            UserDefaults.standard.string(forKey: Settings.keyGraphQLServerURL)
        }
        set {
            UserDefaults.standard.set(newValue, forKey: Settings.keyGraphQLServerURL)
        }
    }
    
    var graphqlServerUsername: String? {
        get {
            UserDefaults.standard.string(forKey: Settings.keyGraphQLServerUsername)
        }
        set {
            UserDefaults.standard.set(newValue, forKey: Settings.keyGraphQLServerUsername)
        }
    }
    
    var graphqlServerPassword: String? {
        get {
            UserDefaults.standard.string(forKey: Settings.keyGraphQLServerPassword)
        }
        set {
            UserDefaults.standard.set(newValue, forKey: Settings.keyGraphQLServerPassword)
        }
    }
    
    func clear() {
        UserDefaults.standard.removeObject(forKey: Settings.keyGraphQLServerURL)
        UserDefaults.standard.removeObject(forKey: Settings.keyGraphQLServerUsername)
        UserDefaults.standard.removeObject(forKey: Settings.keyGraphQLServerPassword)
    }
}
