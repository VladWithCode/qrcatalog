import { queryOptions, mutationOptions } from "@tanstack/react-query";

// Query Keys
export const sectionsQueryKeys = {
    all: () => ["sections"],
    lists: () => [...sectionsQueryKeys.all(), "list"],
    list: (filters: Record<string, any>) => [...sectionsQueryKeys.lists(), filters],
    public: () => [...sectionsQueryKeys.all(), "public"],
    details: () => [...sectionsQueryKeys.all(), "details"],
    detail: (id: string) => [...sectionsQueryKeys.details(), id],
    mutations: () => [...sectionsQueryKeys.all(), "mutations"],
    create: () => [...sectionsQueryKeys.mutations(), "create"],
    update: (id: string) => [...sectionsQueryKeys.mutations(), "update", id],
    delete: (id: string) => [...sectionsQueryKeys.mutations(), "delete", id],
    uploadMedia: (id: string) => [...sectionsQueryKeys.mutations(), "uploadMedia", id],
};

// TypeScript Types
export type TSectionParagraph = {
    id: string;
    section_id: string;
    order: number;
    content: string;
    created_at: string;
    updated_at: string;
};

export type TSectionServiceItem = {
    id: string;
    service_id: string;
    order: number;
    price: number;
    content: string;
    content_as_list: boolean;
    content_list?: string[];
    created_at: string;
    updated_at: string;
};

export type TSectionService = {
    id: string;
    section_id: string;
    title: string;
    price: number;
    description: string;
    items: TSectionServiceItem[];
    created_at: string;
    updated_at: string;
};

export type TSection = {
    id: string;
    name: string;
    title: string;
    image: string;
    bg_image: string;
    paragraphs: TSectionParagraph[];
    services: TSectionService[];
    created_at: string;
    updated_at: string;
};

export type TSectionFilterParams = {
    search?: string;
    has_image?: number;
    has_bg_image?: number;
    paragraph_count?: number;
    service_count?: number;
    item_count?: number;
    min_price?: number;
    max_price?: number;
    created_after?: string;
    created_before?: string;
    updated_after?: string;
    updated_before?: string;
    sort?: string;
    page?: number;
    limit?: number;
};

export type TSectionFilterResult = {
    sections: TSection[];
    total: number;
    page: number;
    limit: number;
    total_pages: number;
    has_next: boolean;
    has_previous: boolean;
    has_error: boolean;
    error: string;
};

export type TCreateSectionData = Omit<TSection, 'id' | 'created_at' | 'updated_at'>;
export type TUpdateSectionData = Partial<TCreateSectionData> & { id: string };
export type TSectionResponse = {
    section: TSection;
    success: boolean;
};
export type TDeleteSectionResponse = {
    success: boolean;
};
export type TUploadMediaData = {
    section_id: string;
    image?: File;
    bg_image?: File;
};
export type TUploadMediaResponse = {
    success: boolean;
    section_id: string;
    image?: string;
    bg_image?: string;
    message: string;
};

// Query Functions
export const getPublicSections = async (): Promise<TSection[]> => {
    const response = await fetch("/api/sections/public", {
        method: "GET",
        credentials: "include",
    });

    if (!response.ok) {
        throw new Error("Failed to fetch public sections");
    }

    const data = await response.json();
    return data.sections || [];
};

export const getSections = async (filters: TSectionFilterParams = {}): Promise<TSectionFilterResult> => {
    const params = new URLSearchParams();

    Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
            params.append(key, value.toString());
        }
    });

    const url = `/api/sections${params.toString() ? `?${params.toString()}` : ''}`;
    const response = await fetch(url, {
        method: "GET",
        credentials: "include",
    });

    if (!response.ok) {
        throw new Error("Failed to fetch sections");
    }

    return await response.json();
};

export const getSectionById = async (sectionId: string): Promise<TSection> => {
    const response = await fetch(`/api/section/${sectionId}`, {
        method: "GET",
        credentials: "include",
    });

    if (!response.ok) {
        throw new Error("Failed to fetch section");
    }

    const data = await response.json();
    return data.section;
};

// Query Options
export const publicSectionsQueryOptions = queryOptions({
    queryKey: sectionsQueryKeys.public(),
    queryFn: getPublicSections,
});

export const sectionsQueryOptions = (filters: TSectionFilterParams = {}) => queryOptions({
    queryKey: sectionsQueryKeys.list(filters),
    queryFn: () => getSections(filters),
});

export const sectionQueryOptions = (id: string) => queryOptions({
    queryKey: sectionsQueryKeys.detail(id),
    queryFn: () => getSectionById(id),
});

// Mutation Functions
export const createSection = async (data: TCreateSectionData): Promise<TSectionResponse> => {
    const response = await fetch("/api/section", {
        method: "POST",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || "Failed to create section");
    }

    return await response.json();
};

export const updateSection = async (data: TUpdateSectionData): Promise<TSectionResponse> => {
    const response = await fetch(`/api/section/${data.id}`, {
        method: "PUT",
        credentials: "include",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || "Failed to update section");
    }

    return await response.json();
};

export const deleteSection = async (id: string): Promise<TDeleteSectionResponse> => {
    const response = await fetch(`/api/section/${id}`, {
        method: "DELETE",
        credentials: "include",
    });

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || "Failed to delete section");
    }

    return await response.json();
};

export const uploadSectionMedia = async (data: TUploadMediaData): Promise<TUploadMediaResponse> => {
    const formData = new FormData();
    formData.append("section_id", data.section_id);

    if (data.image) {
        formData.append("image", data.image);
    }

    if (data.bg_image) {
        formData.append("bg_image", data.bg_image);
    }

    const response = await fetch("/api/sections/media", {
        method: "POST",
        credentials: "include",
        body: formData,
    });

    if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.message || "Failed to upload media");
    }

    return await response.json();
};

// Mutation Options
export const createSectionMutationOptions = mutationOptions({
    mutationKey: sectionsQueryKeys.create(),
    mutationFn: createSection,
});

export const updateSectionMutationOptions = (id: string) => mutationOptions({
    mutationKey: sectionsQueryKeys.update(id),
    mutationFn: updateSection,
});

export const deleteSectionMutationOptions = (id: string) => mutationOptions({
    mutationKey: sectionsQueryKeys.delete(id),
    mutationFn: deleteSection,
});

export const uploadSectionMediaMutationOptions = (id: string) => mutationOptions({
    mutationKey: sectionsQueryKeys.uploadMedia(id),
    mutationFn: uploadSectionMedia,
});
