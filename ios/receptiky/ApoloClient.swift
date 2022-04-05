import Foundation
import Apollo

class Network {
  static let shared = Network()
    
    init() {
        let bundle = Bundle.main.infoDictionary
        
        apollo = ApolloClient(url: URL(string: Bundle.main.object(forInfoDictionaryKey: "API_URL") as! String)!)
    }

    private(set) var apollo: ApolloClient
}
