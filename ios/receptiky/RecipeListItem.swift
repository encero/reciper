//
//  RecipeListItem.swift
//  receptiky
//
//  Created by Matous Michalik on 05.03.2022.
//

import SwiftUI

struct RecipeListItem: View {
    var recipe: Recipe
    var isV1 = false
    
    var body: some View {
        VStack(alignment: .leading){
            if recipe.image != nil {
                HStack{
                    Spacer()
                    Image("noodles")
                        .resizable()
                        .frame( maxWidth: .infinity, maxHeight: 3000)
                        .padding([.leading, .trailing, .top], -28)
                        .scaledToFill()

                    Spacer()
                }
            }
            HStack(alignment: .firstTextBaseline) {
                VStack(alignment: .leading){
                    Text(recipe.title)
                        .font(.title2)
                        .layoutPriority(2)
//                    Text("zluta kucharka")
//                        .font(.caption)
//                        .foregroundColor(.gray)
                }
                Spacer()
                Text("12 dnu")
                    .font(.caption)
                    .frame(minWidth: 30)
            }
        }.clipped()
    }
}

struct RecipeListItem_Previews: PreviewProvider {
    static var previews: some View {
        List {
            RecipeListItem(recipe: recipe, isV1: true)
            RecipeListItem(recipe: recipe, isV1: true)
            RecipeListItem(recipe: recipe)
            RecipeListItem(recipe: Recipe.example)
            RecipeListItem(recipe: Recipe.example, isV1: true)
        }
    }
    
    static var recipe: Recipe {
        var recipe = Recipe.example
        recipe.image = ""
        
        return recipe
    }
}



