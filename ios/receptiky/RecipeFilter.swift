//
//  RecipeFilter.swift
//  receptiky
//
//  Created by Matous Michalik on 19.04.2022.
//

import Foundation

class RecipeFilter {
    enum Sort {
        case alphaAscending
        case alphaDescending
    }
    
    static func filter(_ recipes: [String:Recipe], query: String, sort: Sort, onlyCooked: Bool) -> [Recipe] {
        var list = recipes.values.sorted(by: {
            switch sort {
            case .alphaAscending:
                return $0.title.localizedCompare($1.title) == .orderedAscending
            case .alphaDescending:
                return $0.title.localizedCompare($1.title) == .orderedDescending
            }
        })
        
        if onlyCooked {
            list = list.filter { r in
                return r.cooking
            }
        }
        
        if query != "" {
            let foldedQuery = query.folding(options: .diacriticInsensitive, locale: .current).localizedLowercase
            list = list.filter{ r in
                let foldedTitle = r.title.folding(options: .diacriticInsensitive, locale: .current).localizedLowercase
                
                return foldedTitle.contains(foldedQuery)
            }
        }
        
        return list
    }
}
