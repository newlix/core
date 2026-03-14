// Item is a to-do item.
struct Item: Codable {
    // ID is the unique id
    var id: Int = 0

    // Text is the content
    var text: String = ""

    // CreatedAt is the timestamp which the item created at
    var createdAt: Date = Date()

    enum CodingKeys: String, CodingKey {
        case id = "id"
        case text = "text"
        case createdAt = "created_at"
    }
}

extension Item {
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        if let id = try container.decodeIfPresent(Int.self, forKey: .id) {
            self.id = id
        }

        if let text = try container.decodeIfPresent(String.self, forKey: .text) {
            self.text = text
        }

        if let createdAt = try container.decodeIfPresent(Date.self, forKey: .createdAt) {
            self.createdAt = createdAt
        }

    }
}
