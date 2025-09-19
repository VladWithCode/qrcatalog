import type { QueryClient } from "@tanstack/react-query";
import type { TCheckAuthResponse } from "./auth/service";

export interface IRouterContext {
    auth: TCheckAuthResponse | null;
    queryClient: QueryClient;
}
