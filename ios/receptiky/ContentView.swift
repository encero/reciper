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
    @State private var sort = RecipeFilter.Sort.alphaAscending
    
    @State private var searchTerm = ""
    
    var body: some View {
        NavigationView{
            VStack{
                HStack {
                    cookedFilterSelector()
                    sortSelector()
                    Spacer()
                }
                .padding(.horizontal, 20)
                if case .Error(let error) =  recipes.state {
                    Text("error: " + error)
                }
                List {
                    ForEach(RecipeFilter.filter(
                        recipes.all,
                        query: searchTerm,
                        sort: sort,
                        onlyCooked: filteringCooked
                        
                    ), id: \.id) { recipe in
                        NavigationLink(destination: RecipeDetail(recipe: recipe)) {
                            RecipeListItem(recipe: recipe)
                        }.swipeActions(content: {swipeAction(recipe)})
                            .listRowBackground(highlightCooking(recipe: recipe))
                    }
                }
                .refreshable {
                    await recipes.load()
                }
                .toolbar {
                    ToolbarItem(placement: .navigationBarTrailing) {
                        NavigationLink(destination: RecipeCreate()) {
                            Image(systemName: "plus.circle")
                        }
                    }
                    ToolbarItemGroup(placement: .bottomBar){
                        Spacer()
                        NavigationLink(destination: SettingsView()) {
                            Image(systemName: "gear")
                        }
                        
                    }
                }
                .searchable(text: $searchTerm, placement: .navigationBarDrawer(displayMode: .automatic), prompt: "find recipe")
#if DEBUG
                Button("reset") {
                    Settings.shared.clear()
                    exit(0)
                }
#endif
            }
            .navigationTitle("Receptiky")
        }
        .navigationViewStyle(.stack)
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
    
    func cookedFilterSelector() -> some View {
        Menu {
            Button {
                filteringCooked = false
            } label: {
                Text("vsechny")
            }
            Button {
                filteringCooked = true
            } label: {
                Text("varime")
            }
        } label: {
            HStack {
                Text(filteringCooked ? "varime" : "vsechny")
                Image(systemName: "chevron.down")
                    .font(.system(size: 10))
                    .padding(.leading, -5)
            }
            .padding(.init(top: 5, leading: 10, bottom:5, trailing:10))
            .background(Color(.systemFill))
            .cornerRadius(10)
            .foregroundColor(Color(.secondaryLabel))
            
        }
    }
    
    func sortSelector() -> some View {
        Menu {
            Button {
                sort = .alphaAscending
            } label: {
                sortTypeToView(.alphaAscending)
                
            }
            Button {
                sort = .alphaDescending
            } label: {
                sortTypeToView(.alphaDescending)
            }
        } label: {
            HStack {
                sortTypeToView(sort)
            }
            .padding(.init(top: 5, leading: 10, bottom:5, trailing:10))
            .background(Color(.systemFill))
            .cornerRadius(10)
            .foregroundColor(Color(.secondaryLabel))
            
        }
    }
    
    func sortTypeToView(_ sort: RecipeFilter.Sort) -> some View {
        switch sort {
        case .alphaAscending:
            return HStack {
                Text("Abecedne")
                Image(systemName: "arrow.up")
                    .font(.system(size: 10))
                    .padding(.leading, -5)
            }
        case .alphaDescending:
            return HStack {
                Text("Abecedne")
                Image(systemName: "arrow.down")
                    .font(.system(size: 10))
                    .padding(.leading, -5)
            }
        }
    }
}
