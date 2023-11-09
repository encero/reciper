//
//  GqlTine.swift
//  receptiky
//
//  Created by Matous Michalik on 24.04.2022.
//

import Foundation
import Apollo

public struct Time: JSONDecodable {
    
    let date: Date
    
    public init(jsonValue value: JSONValue) throws {
        let dateFormatter = ISO8601DateFormatter()
        dateFormatter.formatOptions = [.withInternetDateTime, .withFractionalSeconds]
        
        guard let stringValue = value as? String, let date = dateFormatter.date(from: stringValue) else {
                throw JSONDecodingError.couldNotConvert(value: value, to: Time.self)
            }

        self.date = date
    }
}
