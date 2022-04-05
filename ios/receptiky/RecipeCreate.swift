//
//  RecipeCreate.swift
//  receptiky
//
//  Created by Matous Michalik on 06.03.2022.
//

import SwiftUI

struct RecipeCreate: View {
    @EnvironmentObject var recipes:RecipeDataManager
    
    @Environment(\.dismiss) private var dismiss
    
    @State var saving: Bool = false
    @State var title: String = ""
    
    @State private var error: Error?
    
    var body: some View {
        Form{
            Section("Nazev") {
                TextField("Noodles", text: $title)
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
            
        }
        .navigationBarTitleDisplayMode(.inline)
        .navigationTitle("New recipe")
        .navigationBarBackButtonHidden(saving)
        .alert(error?.localizedDescription ?? "", isPresented: .constant(error != nil)) {
            Button("OK") { 
                error = nil
            }
        }
    }
    
    func save() {
        saving = true
        
        recipes.create(title: title) { error in
            saving = false
            
            self.error = error
            if error == nil {
                dismiss()
            }
        }
    }
}

struct RecipeCreate_Previews: PreviewProvider {
    static var previews: some View {
        NavigationView{
            RecipeCreate()
                .environmentObject(RecipeDataManager.example)
        }
    }
}

