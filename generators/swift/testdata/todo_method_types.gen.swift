struct AddItemInput: Codable {
    // the item to add.
    var item: Item = Item()


    enum CodingKeys: String, CodingKey {
        case item = "item"
    }
}

extension AddItemInput {
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        if let item = try container.decodeIfPresent(Item.self, forKey: .item) {
            self.item = item
        }

    }
}

struct AddItemOutput: Codable {
}

struct GetItemsInput: Codable {

}


struct GetItemsOutput: Codable {
    // Items is the list of to-do items.
    var items: [Item] = []

    enum CodingKeys: String, CodingKey {
        case items = "items"
    }
}

extension GetItemsOutput {
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        if let items = try container.decodeIfPresent([Item].self, forKey: .items) {
            self.items = items
        }

    }
}
struct RemoveItemInput: Codable {
    // the id of the item to remove.
    var id: Int = 0


    enum CodingKeys: String, CodingKey {
        case id = "id"
    }
}

extension RemoveItemInput {
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        if let id = try container.decodeIfPresent(Int.self, forKey: .id) {
            self.id = id
        }

    }
}

struct RemoveItemOutput: Codable {
    // the item removed.
    var item: Item = Item()

    enum CodingKeys: String, CodingKey {
        case item = "item"
    }
}

extension RemoveItemOutput {
    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        if let item = try container.decodeIfPresent(Item.self, forKey: .item) {
            self.item = item
        }

    }
}
