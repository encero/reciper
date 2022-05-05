// @generated
//  This file was automatically generated and should not be edited.

import Apollo
import Foundation

public enum Status: RawRepresentable, Equatable, Hashable, CaseIterable, Apollo.JSONDecodable, Apollo.JSONEncodable {
  public typealias RawValue = String
  case success
  case error
  case notFound
  /// Auto generated constant for unknown enum values
  case __unknown(RawValue)

  public init?(rawValue: RawValue) {
    switch rawValue {
      case "Success": self = .success
      case "Error": self = .error
      case "NotFound": self = .notFound
      default: self = .__unknown(rawValue)
    }
  }

  public var rawValue: RawValue {
    switch self {
      case .success: return "Success"
      case .error: return "Error"
      case .notFound: return "NotFound"
      case .__unknown(let value): return value
    }
  }

  public static func == (lhs: Status, rhs: Status) -> Bool {
    switch (lhs, rhs) {
      case (.success, .success): return true
      case (.error, .error): return true
      case (.notFound, .notFound): return true
      case (.__unknown(let lhsValue), .__unknown(let rhsValue)): return lhsValue == rhsValue
      default: return false
    }
  }

  public static var allCases: [Status] {
    return [
      .success,
      .error,
      .notFound,
    ]
  }
}

public final class ListRecipesQuery: GraphQLQuery {
  /// The raw GraphQL definition of this operation.
  public let operationDefinition: String =
    """
    query ListRecipes {
      recipes {
        __typename
        id
        name
        planned
        lastCookedAt
      }
    }
    """

  public let operationName: String = "ListRecipes"

  public init() {
  }

  public struct Data: GraphQLSelectionSet {
    public static let possibleTypes: [String] = ["Query"]

    public static var selections: [GraphQLSelection] {
      return [
        GraphQLField("recipes", type: .nonNull(.list(.nonNull(.object(Recipe.selections))))),
      ]
    }

    public private(set) var resultMap: ResultMap

    public init(unsafeResultMap: ResultMap) {
      self.resultMap = unsafeResultMap
    }

    public init(recipes: [Recipe]) {
      self.init(unsafeResultMap: ["__typename": "Query", "recipes": recipes.map { (value: Recipe) -> ResultMap in value.resultMap }])
    }

    public var recipes: [Recipe] {
      get {
        return (resultMap["recipes"] as! [ResultMap]).map { (value: ResultMap) -> Recipe in Recipe(unsafeResultMap: value) }
      }
      set {
        resultMap.updateValue(newValue.map { (value: Recipe) -> ResultMap in value.resultMap }, forKey: "recipes")
      }
    }

    public struct Recipe: GraphQLSelectionSet {
      public static let possibleTypes: [String] = ["Recipe"]

      public static var selections: [GraphQLSelection] {
        return [
          GraphQLField("__typename", type: .nonNull(.scalar(String.self))),
          GraphQLField("id", type: .nonNull(.scalar(GraphQLID.self))),
          GraphQLField("name", type: .nonNull(.scalar(String.self))),
          GraphQLField("planned", type: .nonNull(.scalar(Bool.self))),
          GraphQLField("lastCookedAt", type: .scalar(Time.self)),
        ]
      }

      public private(set) var resultMap: ResultMap

      public init(unsafeResultMap: ResultMap) {
        self.resultMap = unsafeResultMap
      }

      public init(id: GraphQLID, name: String, planned: Bool, lastCookedAt: Time? = nil) {
        self.init(unsafeResultMap: ["__typename": "Recipe", "id": id, "name": name, "planned": planned, "lastCookedAt": lastCookedAt])
      }

      public var __typename: String {
        get {
          return resultMap["__typename"]! as! String
        }
        set {
          resultMap.updateValue(newValue, forKey: "__typename")
        }
      }

      public var id: GraphQLID {
        get {
          return resultMap["id"]! as! GraphQLID
        }
        set {
          resultMap.updateValue(newValue, forKey: "id")
        }
      }

      public var name: String {
        get {
          return resultMap["name"]! as! String
        }
        set {
          resultMap.updateValue(newValue, forKey: "name")
        }
      }

      public var planned: Bool {
        get {
          return resultMap["planned"]! as! Bool
        }
        set {
          resultMap.updateValue(newValue, forKey: "planned")
        }
      }

      public var lastCookedAt: Time? {
        get {
          return resultMap["lastCookedAt"] as? Time
        }
        set {
          resultMap.updateValue(newValue, forKey: "lastCookedAt")
        }
      }
    }
  }
}

public final class UpdateRecipeMutation: GraphQLMutation {
  /// The raw GraphQL definition of this operation.
  public let operationDefinition: String =
    """
    mutation UpdateRecipe($id: ID!, $name: String!) {
      updateRecipe(input: {id: $id, name: $name}) {
        __typename
        status
      }
    }
    """

  public let operationName: String = "UpdateRecipe"

  public var id: GraphQLID
  public var name: String

  public init(id: GraphQLID, name: String) {
    self.id = id
    self.name = name
  }

  public var variables: GraphQLMap? {
    return ["id": id, "name": name]
  }

  public struct Data: GraphQLSelectionSet {
    public static let possibleTypes: [String] = ["Mutation"]

    public static var selections: [GraphQLSelection] {
      return [
        GraphQLField("updateRecipe", arguments: ["input": ["id": GraphQLVariable("id"), "name": GraphQLVariable("name")]], type: .nonNull(.object(UpdateRecipe.selections))),
      ]
    }

    public private(set) var resultMap: ResultMap

    public init(unsafeResultMap: ResultMap) {
      self.resultMap = unsafeResultMap
    }

    public init(updateRecipe: UpdateRecipe) {
      self.init(unsafeResultMap: ["__typename": "Mutation", "updateRecipe": updateRecipe.resultMap])
    }

    public var updateRecipe: UpdateRecipe {
      get {
        return UpdateRecipe(unsafeResultMap: resultMap["updateRecipe"]! as! ResultMap)
      }
      set {
        resultMap.updateValue(newValue.resultMap, forKey: "updateRecipe")
      }
    }

    public struct UpdateRecipe: GraphQLSelectionSet {
      public static let possibleTypes: [String] = ["Result"]

      public static var selections: [GraphQLSelection] {
        return [
          GraphQLField("__typename", type: .nonNull(.scalar(String.self))),
          GraphQLField("status", type: .nonNull(.scalar(Status.self))),
        ]
      }

      public private(set) var resultMap: ResultMap

      public init(unsafeResultMap: ResultMap) {
        self.resultMap = unsafeResultMap
      }

      public init(status: Status) {
        self.init(unsafeResultMap: ["__typename": "Result", "status": status])
      }

      public var __typename: String {
        get {
          return resultMap["__typename"]! as! String
        }
        set {
          resultMap.updateValue(newValue, forKey: "__typename")
        }
      }

      public var status: Status {
        get {
          return resultMap["status"]! as! Status
        }
        set {
          resultMap.updateValue(newValue, forKey: "status")
        }
      }
    }
  }
}

public final class PlanRecipeMutation: GraphQLMutation {
  /// The raw GraphQL definition of this operation.
  public let operationDefinition: String =
    """
    mutation PlanRecipe($id: ID!) {
      planRecipe(id: $id) {
        __typename
        status
      }
    }
    """

  public let operationName: String = "PlanRecipe"

  public var id: GraphQLID

  public init(id: GraphQLID) {
    self.id = id
  }

  public var variables: GraphQLMap? {
    return ["id": id]
  }

  public struct Data: GraphQLSelectionSet {
    public static let possibleTypes: [String] = ["Mutation"]

    public static var selections: [GraphQLSelection] {
      return [
        GraphQLField("planRecipe", arguments: ["id": GraphQLVariable("id")], type: .nonNull(.object(PlanRecipe.selections))),
      ]
    }

    public private(set) var resultMap: ResultMap

    public init(unsafeResultMap: ResultMap) {
      self.resultMap = unsafeResultMap
    }

    public init(planRecipe: PlanRecipe) {
      self.init(unsafeResultMap: ["__typename": "Mutation", "planRecipe": planRecipe.resultMap])
    }

    public var planRecipe: PlanRecipe {
      get {
        return PlanRecipe(unsafeResultMap: resultMap["planRecipe"]! as! ResultMap)
      }
      set {
        resultMap.updateValue(newValue.resultMap, forKey: "planRecipe")
      }
    }

    public struct PlanRecipe: GraphQLSelectionSet {
      public static let possibleTypes: [String] = ["Result"]

      public static var selections: [GraphQLSelection] {
        return [
          GraphQLField("__typename", type: .nonNull(.scalar(String.self))),
          GraphQLField("status", type: .nonNull(.scalar(Status.self))),
        ]
      }

      public private(set) var resultMap: ResultMap

      public init(unsafeResultMap: ResultMap) {
        self.resultMap = unsafeResultMap
      }

      public init(status: Status) {
        self.init(unsafeResultMap: ["__typename": "Result", "status": status])
      }

      public var __typename: String {
        get {
          return resultMap["__typename"]! as! String
        }
        set {
          resultMap.updateValue(newValue, forKey: "__typename")
        }
      }

      public var status: Status {
        get {
          return resultMap["status"]! as! Status
        }
        set {
          resultMap.updateValue(newValue, forKey: "status")
        }
      }
    }
  }
}

public final class UnPlanRecipeMutation: GraphQLMutation {
  /// The raw GraphQL definition of this operation.
  public let operationDefinition: String =
    """
    mutation UnPlanRecipe($id: ID!) {
      unPlanRecipe(id: $id) {
        __typename
        status
      }
    }
    """

  public let operationName: String = "UnPlanRecipe"

  public var id: GraphQLID

  public init(id: GraphQLID) {
    self.id = id
  }

  public var variables: GraphQLMap? {
    return ["id": id]
  }

  public struct Data: GraphQLSelectionSet {
    public static let possibleTypes: [String] = ["Mutation"]

    public static var selections: [GraphQLSelection] {
      return [
        GraphQLField("unPlanRecipe", arguments: ["id": GraphQLVariable("id")], type: .nonNull(.object(UnPlanRecipe.selections))),
      ]
    }

    public private(set) var resultMap: ResultMap

    public init(unsafeResultMap: ResultMap) {
      self.resultMap = unsafeResultMap
    }

    public init(unPlanRecipe: UnPlanRecipe) {
      self.init(unsafeResultMap: ["__typename": "Mutation", "unPlanRecipe": unPlanRecipe.resultMap])
    }

    public var unPlanRecipe: UnPlanRecipe {
      get {
        return UnPlanRecipe(unsafeResultMap: resultMap["unPlanRecipe"]! as! ResultMap)
      }
      set {
        resultMap.updateValue(newValue.resultMap, forKey: "unPlanRecipe")
      }
    }

    public struct UnPlanRecipe: GraphQLSelectionSet {
      public static let possibleTypes: [String] = ["Result"]

      public static var selections: [GraphQLSelection] {
        return [
          GraphQLField("__typename", type: .nonNull(.scalar(String.self))),
          GraphQLField("status", type: .nonNull(.scalar(Status.self))),
        ]
      }

      public private(set) var resultMap: ResultMap

      public init(unsafeResultMap: ResultMap) {
        self.resultMap = unsafeResultMap
      }

      public init(status: Status) {
        self.init(unsafeResultMap: ["__typename": "Result", "status": status])
      }

      public var __typename: String {
        get {
          return resultMap["__typename"]! as! String
        }
        set {
          resultMap.updateValue(newValue, forKey: "__typename")
        }
      }

      public var status: Status {
        get {
          return resultMap["status"]! as! Status
        }
        set {
          resultMap.updateValue(newValue, forKey: "status")
        }
      }
    }
  }
}

public final class CookRecipeMutation: GraphQLMutation {
  /// The raw GraphQL definition of this operation.
  public let operationDefinition: String =
    """
    mutation CookRecipe($id: ID!) {
      cookRecipe(id: $id) {
        __typename
        status
      }
    }
    """

  public let operationName: String = "CookRecipe"

  public var id: GraphQLID

  public init(id: GraphQLID) {
    self.id = id
  }

  public var variables: GraphQLMap? {
    return ["id": id]
  }

  public struct Data: GraphQLSelectionSet {
    public static let possibleTypes: [String] = ["Mutation"]

    public static var selections: [GraphQLSelection] {
      return [
        GraphQLField("cookRecipe", arguments: ["id": GraphQLVariable("id")], type: .nonNull(.object(CookRecipe.selections))),
      ]
    }

    public private(set) var resultMap: ResultMap

    public init(unsafeResultMap: ResultMap) {
      self.resultMap = unsafeResultMap
    }

    public init(cookRecipe: CookRecipe) {
      self.init(unsafeResultMap: ["__typename": "Mutation", "cookRecipe": cookRecipe.resultMap])
    }

    public var cookRecipe: CookRecipe {
      get {
        return CookRecipe(unsafeResultMap: resultMap["cookRecipe"]! as! ResultMap)
      }
      set {
        resultMap.updateValue(newValue.resultMap, forKey: "cookRecipe")
      }
    }

    public struct CookRecipe: GraphQLSelectionSet {
      public static let possibleTypes: [String] = ["Result"]

      public static var selections: [GraphQLSelection] {
        return [
          GraphQLField("__typename", type: .nonNull(.scalar(String.self))),
          GraphQLField("status", type: .nonNull(.scalar(Status.self))),
        ]
      }

      public private(set) var resultMap: ResultMap

      public init(unsafeResultMap: ResultMap) {
        self.resultMap = unsafeResultMap
      }

      public init(status: Status) {
        self.init(unsafeResultMap: ["__typename": "Result", "status": status])
      }

      public var __typename: String {
        get {
          return resultMap["__typename"]! as! String
        }
        set {
          resultMap.updateValue(newValue, forKey: "__typename")
        }
      }

      public var status: Status {
        get {
          return resultMap["status"]! as! Status
        }
        set {
          resultMap.updateValue(newValue, forKey: "status")
        }
      }
    }
  }
}

public final class CreateRecipeMutation: GraphQLMutation {
  /// The raw GraphQL definition of this operation.
  public let operationDefinition: String =
    """
    mutation CreateRecipe($name: String!) {
      createRecipe(input: {name: $name}) {
        __typename
        id
        name
      }
    }
    """

  public let operationName: String = "CreateRecipe"

  public var name: String

  public init(name: String) {
    self.name = name
  }

  public var variables: GraphQLMap? {
    return ["name": name]
  }

  public struct Data: GraphQLSelectionSet {
    public static let possibleTypes: [String] = ["Mutation"]

    public static var selections: [GraphQLSelection] {
      return [
        GraphQLField("createRecipe", arguments: ["input": ["name": GraphQLVariable("name")]], type: .nonNull(.object(CreateRecipe.selections))),
      ]
    }

    public private(set) var resultMap: ResultMap

    public init(unsafeResultMap: ResultMap) {
      self.resultMap = unsafeResultMap
    }

    public init(createRecipe: CreateRecipe) {
      self.init(unsafeResultMap: ["__typename": "Mutation", "createRecipe": createRecipe.resultMap])
    }

    public var createRecipe: CreateRecipe {
      get {
        return CreateRecipe(unsafeResultMap: resultMap["createRecipe"]! as! ResultMap)
      }
      set {
        resultMap.updateValue(newValue.resultMap, forKey: "createRecipe")
      }
    }

    public struct CreateRecipe: GraphQLSelectionSet {
      public static let possibleTypes: [String] = ["Recipe"]

      public static var selections: [GraphQLSelection] {
        return [
          GraphQLField("__typename", type: .nonNull(.scalar(String.self))),
          GraphQLField("id", type: .nonNull(.scalar(GraphQLID.self))),
          GraphQLField("name", type: .nonNull(.scalar(String.self))),
        ]
      }

      public private(set) var resultMap: ResultMap

      public init(unsafeResultMap: ResultMap) {
        self.resultMap = unsafeResultMap
      }

      public init(id: GraphQLID, name: String) {
        self.init(unsafeResultMap: ["__typename": "Recipe", "id": id, "name": name])
      }

      public var __typename: String {
        get {
          return resultMap["__typename"]! as! String
        }
        set {
          resultMap.updateValue(newValue, forKey: "__typename")
        }
      }

      public var id: GraphQLID {
        get {
          return resultMap["id"]! as! GraphQLID
        }
        set {
          resultMap.updateValue(newValue, forKey: "id")
        }
      }

      public var name: String {
        get {
          return resultMap["name"]! as! String
        }
        set {
          resultMap.updateValue(newValue, forKey: "name")
        }
      }
    }
  }
}

public final class DeleteRecipeMutation: GraphQLMutation {
  /// The raw GraphQL definition of this operation.
  public let operationDefinition: String =
    """
    mutation DeleteRecipe($id: ID!) {
      deleteRecipe(id: $id) {
        __typename
        status
      }
    }
    """

  public let operationName: String = "DeleteRecipe"

  public var id: GraphQLID

  public init(id: GraphQLID) {
    self.id = id
  }

  public var variables: GraphQLMap? {
    return ["id": id]
  }

  public struct Data: GraphQLSelectionSet {
    public static let possibleTypes: [String] = ["Mutation"]

    public static var selections: [GraphQLSelection] {
      return [
        GraphQLField("deleteRecipe", arguments: ["id": GraphQLVariable("id")], type: .nonNull(.object(DeleteRecipe.selections))),
      ]
    }

    public private(set) var resultMap: ResultMap

    public init(unsafeResultMap: ResultMap) {
      self.resultMap = unsafeResultMap
    }

    public init(deleteRecipe: DeleteRecipe) {
      self.init(unsafeResultMap: ["__typename": "Mutation", "deleteRecipe": deleteRecipe.resultMap])
    }

    public var deleteRecipe: DeleteRecipe {
      get {
        return DeleteRecipe(unsafeResultMap: resultMap["deleteRecipe"]! as! ResultMap)
      }
      set {
        resultMap.updateValue(newValue.resultMap, forKey: "deleteRecipe")
      }
    }

    public struct DeleteRecipe: GraphQLSelectionSet {
      public static let possibleTypes: [String] = ["Result"]

      public static var selections: [GraphQLSelection] {
        return [
          GraphQLField("__typename", type: .nonNull(.scalar(String.self))),
          GraphQLField("status", type: .nonNull(.scalar(Status.self))),
        ]
      }

      public private(set) var resultMap: ResultMap

      public init(unsafeResultMap: ResultMap) {
        self.resultMap = unsafeResultMap
      }

      public init(status: Status) {
        self.init(unsafeResultMap: ["__typename": "Result", "status": status])
      }

      public var __typename: String {
        get {
          return resultMap["__typename"]! as! String
        }
        set {
          resultMap.updateValue(newValue, forKey: "__typename")
        }
      }

      public var status: Status {
        get {
          return resultMap["status"]! as! Status
        }
        set {
          resultMap.updateValue(newValue, forKey: "status")
        }
      }
    }
  }
}
