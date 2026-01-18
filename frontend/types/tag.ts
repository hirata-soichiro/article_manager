export interface Tag {
    id: number
    name: string
    createdAt: string
    updatedAt: string
}

export interface CreateTagInput {
    name: string
}

export interface UpdateTagInput {
    name: string
}
