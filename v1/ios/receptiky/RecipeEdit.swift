//
//  RecipeEdit.swift
//  receptiky
//
//  Created by Matous Michalik on 06.03.2022.
//

import SwiftUI

struct RecipeEdit: View {
    @EnvironmentObject var recipes:RecipeDataManager
    
    @Environment(\.dismiss) private var dismiss
    
    @State var recipe: Recipe
    
    @State private var saving: Bool = false
    @State private var deleting: Bool = false
    
    @State private var error: Error?
    
    @State private var presentDeleteConfirmation: Bool = false
    
    var body: some View {
        Form{
            Section("Nazev") {
                TextField("Noodles", text: $recipe.title)
            }
            //            Section("Zdroj") {
            //                TextField("kucharka", text: $recipe.source)
            //            }
            Button(action: save) {
                HStack{
                    Text("Ulozit")
                    if saving {
                        Spacer()
                        ProgressView()
                            .progressViewStyle(.circular)
                    }
                }
            }
            Button(action: confirmDelete) {
                HStack{
                    Text("Delete")
                        .foregroundColor(.red)
                    if deleting {
                        Spacer()
                        ProgressView()
                            .progressViewStyle(.circular)
                    }
                }
            }
            
        }
        .navigationBarTitleDisplayMode(.inline)
        .navigationTitle(recipe.title)
        .navigationBarBackButtonHidden(saving)
        .alert(error?.localizedDescription ?? "", isPresented: .constant(error != nil)) {
            Button("OK", role:.cancel) {
                error = nil
            }
        }
        .confirmationDialog("Confirm deletion?", isPresented: $presentDeleteConfirmation) {
            Button("Yes", role: .destructive, action: performDelete)
        }
    }
    
    func performDelete() {
        recipes.delete(recipe) { error in
            deleting = false
            
            self.error = error
            if error == nil {
                dismiss()
            }
        }
    }
    
    func confirmDelete() {
        guard !deleting && !saving else { return }
        
        presentDeleteConfirmation = true
    }
    
    func save() {
        guard !deleting && !saving else { return }
        
        saving = true
        recipes.updateRecipe(recipe) { error in
            saving = false
            
            self.error = error
            if error == nil {
                dismiss()
            }
        }
    }
}

struct RecipeEdit_Previews: PreviewProvider {
    static var previews: some View {
        NavigationView{
            RecipeEdit(recipe: Recipe.example)
                .environmentObject(RecipeDataManager.example)
        }
    }
}

