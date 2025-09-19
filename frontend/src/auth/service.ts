import { queryOptions } from "@tanstack/react-query";

export const authQueryKeys = {
    all: () => ["auth"],
    isAuthenticated: () => [...authQueryKeys.all(), "isAuthenticated"],
    user: () => [...authQueryKeys.all(), "user"],
};

export type TCheckAuthResponse = {
    isAuthenticated: boolean;
    user?: string;
};

export const checkAuth = async () => {
    const response = await fetch("/api/auth", {
        method: "GET",
        credentials: "include",
        mode: "cors",
    });

    if (!response.ok) {
        throw new Error("Error al verificar autenticación");
    }

    return response.json() as Promise<TCheckAuthResponse>;
};

export const checkAuthOptions = queryOptions({
    queryKey: authQueryKeys.isAuthenticated(),
    queryFn: checkAuth,
});

