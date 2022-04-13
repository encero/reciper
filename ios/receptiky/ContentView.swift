//
//  ContentView.swift
//  receptiky
//
//  Created by Matous Michalik on 05.03.2022.
//

import SwiftUI

struct ContentView: View {
    @EnvironmentObject var recipes: RecipeDataManager
    
    @State private var filteringCooked = false
    
    var body: some View {
        NavigationView{
            VStack{
                List {
                    if case .Error(let error) =  recipes.state {
                        Text("error: " + error)
                    }
                    ForEach((filteringCooked ? recipes.onlyCooked : recipes.all).sorted(by: {$0.value.title < $1.value.title}), id: \.key) { _,recipe in
                        NavigationLink(destination: recipeDetail(recipe: recipe)) {
                            RecipeListItem(recipe: recipe)
                        }.swipeActions(content: {swipeAction(recipe)})
                            .listRowBackground(highlightCooking(recipe: recipe))
                    }
                }
                .refreshable {
                    await recipes.load()
                }
                .navigationTitle("Receptiky")
                .toolbar {
                    ToolbarItem(placement: .navigationBarTrailing) {
                        NavigationLink(destination: RecipeCreate()) {
                            Image(systemName: "plus.circle")
                        }
                    }
                    ToolbarItemGroup(placement: .bottomBar){
                        Button(filteringCooked ? "vsechny" : "varime") {
                            withAnimation{
                                filteringCooked.toggle()
                            }
                        }
                        NavigationLink(destination: SettingsView()) {
                            Image(systemName: "gear")
                        }
                        
                    }
                }
                #if DEBUG
                Button("reset") {
                    Settings.shared.clear()
                    exit(0)
                }
                #endif
            }
        }.navigationViewStyle(.stack)
    }
    
    @ViewBuilder func swipeAction(_ recipe:Recipe) -> some View {
        Button(recipe.cooking ? "Uvareno" : "Uvarit") {
            if !recipe.cooking {
                recipes.plan(recipe)
            } else {
                recipes.cooking(recipe)
            }
        }.tint(recipe.cooking ? .yellow : .green)

        if recipe.cooking {
            Button("Zrusit") {
                recipes.unPlan(recipe)
            }
        }
        
    }
    
    func highlightCooking(recipe: Recipe) -> Color {
        return recipe.cooking ? .yellow.opacity(0.3) : Color(.systemBackground)
    }
    
    func recipeDetail(recipe: Recipe) -> some View {
        return RecipeDetail(recipe: recipe)
    }
}
