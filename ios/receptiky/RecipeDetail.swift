//
//  RecipeDetail.swift
//  receptiky
//
//  Created by Matous Michalik on 06.03.2022.
//

import SwiftUI

struct RecipeDetail: View {
    var recipe:Recipe
    
    @EnvironmentObject var recipes:RecipeDataManager
    
    @Environment(\.dismiss) private var dismiss
    
    @State private var error: Error?
    @State private var updating: Bool = false
    
    var body: some View {
        VStack{
            Image("noodles")
                .resizable()
                .scaledToFit()
                .frame(maxHeight:200)
            Button(action: self.toggleCooking) {
                Text(recipe.cooking ? "Uvareno": "Uvarit")
                    .frame(maxWidth:.infinity)
                
                if updating {
                    Spacer()
                    ProgressView()
                        .progressViewStyle(.circular)
                }
            }
            .padding()
            .background(recipe.cooking ? .yellow: .green)
            .foregroundColor(.white)
            Spacer()
        }
        .navigationBarTitleDisplayMode(.inline)
        .navigationTitle(recipe.title)
        .toolbar {
            NavigationLink(destination: RecipeEdit(recipe:recipe)) {
                Text("upravit")
            }.isDetailLink(false)
        }
        .alert(error?.localizedDescription ?? "", isPresented: .constant(error != nil)) {
            Button("OK") {
                error = nil
            }
        }
    }
    
    func toggleCooking() {
        updating = true
        
        if !recipe.cooking {
            recipes.plan(recipe) { error in
                updating = false
                self.error = error
            }
        } else {
            recipes.cooking(recipe) { error in
                updating = false
                self.error = error
            }
        }
    }
}

struct RecipeDetail_Previews: PreviewProvider {
    static var previews: some View {
        NavigationView{
            RecipeDetail(
                recipe: Recipe.example
            )
        }.environmentObject(RecipeDataManager.example).previewInterfaceOrientation(.portraitUpsideDown)
    }
}

