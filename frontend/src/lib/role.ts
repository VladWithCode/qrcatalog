export type TRole = "admin" | "editor" | "user";

export const hasAccess = (role: TRole, requiredRole: TRole): boolean => {
    switch (requiredRole) {
        case "admin":
            return role === "admin";
        case "editor":
            return role === "admin" || role === "editor";
        case "user":
            return role === "admin" || role === "editor" || role === "user";
        default:
            return false;
    }
};
