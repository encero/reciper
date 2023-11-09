import Foundation
import Apollo

class Network {
    static let shared = Network()
    
    private(set) lazy var apollo: ApolloClient = {
        let store = ApolloStore(cache: InMemoryNormalizedCache())
        let provider = NetworkInterceptorProvider(store: store)
        
        var url = URL(string: Settings.shared.graphqlServerURL ?? "") ?? dummyURL()
        url.appendPathComponent("query")

        let auth = Settings.shared.graphqlServerPassword
        if auth != nil {
            provider.addInterceptor(BasicAuthInterceptor(username: "auth", password: auth!))
        }
        
        let transport = RequestChainNetworkTransport(interceptorProvider: provider,
                                                     endpointURL: url)
        return ApolloClient(networkTransport: transport, store: store)
    }()
    
    private func dummyURL() -> URL {
        return URL(string: "https://127.0.0.1")!
    }
}


class NetworkInterceptorProvider: DefaultInterceptorProvider {
    
    private var addionalInterceptors: [ApolloInterceptor] = []
    
    func addInterceptor(_ interceptor: ApolloInterceptor) {
        addionalInterceptors.append(interceptor)
    }
    
    override func interceptors<Operation: GraphQLOperation>(for operation: Operation) -> [ApolloInterceptor] {
        var interceptors = super.interceptors(for: operation)
        
        for interceptor in addionalInterceptors {
            interceptors.insert(interceptor, at: 0)
        }
        
        return interceptors
    }
}

class BasicAuthInterceptor: ApolloInterceptor {
    
    init(username:String, password:String) {
        headerValue = Data("\(username):\(password)".utf8).base64EncodedString()
    }
    
    private var headerValue: String
    
    func interceptAsync<Operation: GraphQLOperation>(
        chain: RequestChain,
        request: HTTPRequest<Operation>,
        response: HTTPResponse<Operation>?,
        completion: @escaping (Swift.Result<GraphQLResult<Operation.Data>, Error>) -> Void) {
            request.addHeader(name: "Authorization", value: "Basic \(headerValue)")
            
            chain.proceedAsync(request: request,
                               response: response,
                               completion: completion)
        }
}
