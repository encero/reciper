//
//  RecipeDataManager.swift
//  receptiky
//
//  Created by Matous Michalik on 05.03.2022.
//

import Foundation

struct Recipe: Identifiable {
    var id: String
    var title: String
    var source: String = ""
    var image: String?
    var lastCookedAt: Date?
    
    var cooking: Bool = false
    
    var lastCookedSince: String {
        guard lastCookedAt != nil else {
            return "nikdy"
        }
        
        guard let days = Calendar.current.dateComponents([.day], from: lastCookedAt!, to: Date()).day else {
            return "nikdy"
        }
        
        if days == 0 {
            return "dnes"
        }
        
        if days == 1 {
            return "vcera"
        }
        
        if days < 7 {
            return "pred \(days) dny"
        }
        
        if days < 30 {
            return "pred mesicem"
        }
        
        return "davno"
    }
    
    static var example: Recipe {
        Recipe(id: "id", title: "Nudlovy salat", cooking: false)
    }
}

class RecipeDataManager: ObservableObject {
    enum State {
        case Loaded
        case Pending
        case Error(String)
    }
    
    @Published var state = State.Pending
    @Published var all = [String:Recipe]()
    
    init() {
        Task {
            await load()
        }
    }
    
    func load() async {
        print("loading recipes")
        Network.shared.apollo.fetch(query: ListRecipesQuery(), cachePolicy: .fetchIgnoringCacheCompletely) { result in
            self.state = State.Loaded
            
            switch result {
            case .success(let graphQLResult):
                print("loading success")
                
                if graphQLResult.data!.recipes.isEmpty {
                    return
                }
                
                let recipes = graphQLResult.data!.recipes
                
                for i in (0...recipes.count - 1) {
                    
                    let r = recipes[i]
                    print("loaded #\(r.id):\(r.name)")
                    
                    let recipe = Recipe(
                        id: r.id,
                        title: r.name,
                        source: "",
                        image: nil,
                        lastCookedAt: r.lastCookedAt?.date,
                        cooking: r.planned
                    )
                    
                    self.all[recipe.id] = recipe
                }
                
            case .failure(let error):
                self.state = .Error(error.localizedDescription)
            }
        }
    }
    
    var onlyCooked: [String:Recipe] {
        all.filter{ _, r in return r.cooking }
    }
    
    func create(title: String, done:((Error?)->())? = nil) {
        let mutation = CreateRecipeMutation(name: title)
        
        Network.shared.apollo.perform(mutation: mutation) { [weak self] result in
            var error: Error?
            defer { done?(error) }
            
            switch result {
            case .success(let success):
                guard let data = success.data?.createRecipe else {return}
                
                self?.all[data.id] = Recipe(id: data.id, title: data.name)
            case .failure(let err):
                print("create recipe failed")
                error = err
            }
        }
    }
    
    func delete(_ recipe:Recipe, done:((Error?)->())? = nil) {
        guard all[recipe.id] != nil else { done?(nil); return }
        
        print("deleting recipe \(recipe.id):\(recipe.title)")
        
        let mutation = DeleteRecipeMutation(id: recipe.id)
        
        Network.shared.apollo.perform(mutation: mutation) { [weak self] result in
            var error: Error?
            defer { done?(error) }
            
            switch result {
            case .success(let result):
                print("delete recipe status \(result.data?.deleteRecipe.status.rawValue ?? "none")")
                
                self?.all[recipe.id] =  nil
            case .failure(let err):
                print("update recipe failure")
                error = err
            }
        }
    }
    
    func updateRecipe(_ recipe:Recipe, done: ((Error?)->())? = nil) {
        guard all[recipe.id] != nil else { done?(nil); return }
        
        print("updating recipe \(recipe)")
        
        let mutation = UpdateRecipeMutation(id: recipe.id, name: recipe.title)
        
        Network.shared.apollo.perform(mutation: mutation) { [weak self] result in
            var error: Error?
            defer { done?(error) }
            
            switch result {
            case .success(let result):
                print("update recipe status \(result.data?.updateRecipe.status.rawValue ?? "none")")
                
                self?.all[recipe.id]!.title = recipe.title
                self?.all[recipe.id]!.cooking = recipe.cooking
            case .failure(let err):
                print("update recipe failure")
                error = err
            }
        }
    }
    
    func plan(_ recipe:Recipe, done:((Error?)->())? = nil) {
        guard all[recipe.id] != nil else { done?(nil); return }
        
        print("planing recipe \(recipe)")
        
        let mutation = PlanRecipeMutation(id: recipe.id)
        
        Network.shared.apollo.perform(mutation: mutation) { [weak self] result in
            var error: Error?
            defer { done?(error) }
            
            switch result {
            case .success(let result):
                print("update recipe status \(result.data?.planRecipe.status.rawValue ?? "none")")
                
                self?.all[recipe.id]!.cooking = true
            case .failure(let err):
                print("update recipe failure")
                error = err
            }
        }
    }
    
    func unPlan(_ recipe:Recipe, done: ((Error?)->())? = nil) {
        guard all[recipe.id] != nil else { done?(nil);return }
        
        print("UN planning recipe \(recipe.id)")
        
        let mutation = UnPlanRecipeMutation(id: recipe.id)
        
        Network.shared.apollo.perform(mutation: mutation) { [weak self] result in
            var error: Error?
            defer { done?(error) }
            
            switch result {
            case .success(let result):
                print("un plan status \(result.data?.unPlanRecipe.status.rawValue ?? "none")")
                
                self?.all[recipe.id]!.cooking = false
            case .failure(let err):
                print("update recipe failure")
                error = err
            }
        }
    }
    
    func cooking(_ recipe:Recipe, done:((Error?)->())? = nil) {
        guard all[recipe.id] != nil else { done?(nil); return }
        
        print("set recipe as cooked \(recipe)")
        
        let mutation = CookRecipeMutation(id: recipe.id)
        
        Network.shared.apollo.perform(mutation: mutation) { [weak self] result in
            var error: Error?
            defer { done?(error) }
            
            switch result {
            case .success(let result):
                print("update recipe status \(result.data?.cookRecipe.status.rawValue ?? "none")")
                
                self?.all[recipe.id]!.cooking = false
            case .failure(let err):
                print("update recipe failure")
                error = err
            }
        }
    }
    
    static var example: RecipeDataManager {
        let manager = RecipeDataManager()
        
        manager.all = [
            "id 1": Recipe(id: "id 1", title: "recipe 1", source: "source", image: nil, cooking: false),
            "id 2": Recipe(id: "id 2", title: "recipe 2", source: "source", image: nil, cooking: true),
        ]
        
        return manager
    }
}
