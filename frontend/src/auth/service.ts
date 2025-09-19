import { mutationOptions, queryOptions } from "@tanstack/react-query";

export const authQueryKeys = {
    all: () => ["auth"],
    isAuthenticated: () => [...authQueryKeys.all(), "isAuthenticated"],
    user: () => [...authQueryKeys.all(), "user"],
    signIn: () => [...authQueryKeys.all(), "signIn"],
};

export type TCheckAuthResponse = {
    isAuthenticated: boolean;
    user?: string;
};

export const checkAuth = async () => {
    try {
        const response = await fetch("/api/auth", {
            method: "GET",
            credentials: "include",
        });
        return await response.json() as TCheckAuthResponse;
    } catch (error) {
        console.error("Error al verificar autenticación:", error);
        return { isAuthenticated: false, user: null };
    }
};

export const checkAuthOptions = queryOptions({
    queryKey: authQueryKeys.isAuthenticated(),
    queryFn: checkAuth,
});

export type TSignInData = {
    username: string;
    password: string;
};

export type TSignInResponse = {
    isAuthenticated: boolean;
    user?: string;
    message?: string;
    redirect?: string;
};

export const signIn = async (data: TSignInData) => {
    let response: Response;

    try {
        response = await fetch("/api/sign-in", {
            method: "POST",
            credentials: "include",
            mode: "cors",
            body: JSON.stringify(data),
        });
    } catch (error) {
        console.error("Error al autenticar usuario:", error);
        throw new Error("Error al autenticar usuario");
    }
    console.log("response", response);

    if (!response.ok) {
        return null;
    }

    return response.json() as Promise<TSignInResponse>;
};

export const signInOptions = mutationOptions({
    mutationKey: authQueryKeys.signIn(),
    mutationFn: signIn,
})
