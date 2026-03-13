struct CoreError: LocalizedError {
	let status: Int
	let message: String
	
	var errorDescription: String? {
		message
	}
}

// TodoClient is the API client.
struct TodoClient {
	// encoder is the conventional json encoder
	private let encoder = JSONEncoder()
	
	// decoder is the conventional json decoder
	private let decoder = JSONDecoder()
	
	// endpoint is the required API endpoint address.
	let endpoint: String

	// AuthToken is an optional authentication token.
	var authToken: String?

	// session is the client used for making requests, defaulting to URLSession.shared.
	let session: URLSession = URLSession.shared

	private func call<Input, Output>(method: String, input: Input) async throws -> Output where Input: Codable, Output: Codable {
		guard let url = URL(string: endpoint + "/" + method) else {
			throw CoreError(status: 0, message: "Invalid URL: \(endpoint)/\(method)")
		}

		let body = try self.encoder.encode(input)

		var req = URLRequest(url: url)
		req.setValue("application/json", forHTTPHeaderField: "Content-Type")
		if let tok = self.authToken {
			req.setValue("Bearer " + tok, forHTTPHeaderField: "Authorization")
		}
		req.httpMethod = "POST"
		req.httpBody = body

		let (data, res) = try await self.session.data(for: req)

		guard let r = res as? HTTPURLResponse else {
			throw CoreError(status: 0, message: "Unexpected response type")
		}
		
		if r.statusCode >= 300 {
			let body = String(decoding: data, as: UTF8.self)
			let err = CoreError(status: r.statusCode, message: body)
			throw err
		}
		return try self.decoder.decode(Output.self, from: data)
	}


	// AddItem adds an item to the list.
	func addItem(input: AddItemInput) async throws -> AddItemOutput {
		return try await call(method: "add_item", input: input)
	}

	// GetItems returns all items in the list.
	func getItems(input: GetItemsInput) async throws -> GetItemsOutput {
		return try await call(method: "get_items", input: input)
	}

	// RemoveItem removes an item from the to-do list.
	func removeItem(input: RemoveItemInput) async throws -> RemoveItemOutput {
		return try await call(method: "remove_item", input: input)
	}

}

