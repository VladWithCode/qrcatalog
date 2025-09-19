import { Footer } from "@/components/footer";
import { Header } from "@/components/header";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/iniciar-sesion")({
    component: RouteComponent,
    staticData: {
        withOpaqueHeader: true,
    },
});

function RouteComponent() {
    return (
        <>
            <Header noAnimate={true} alwaysOpaque={true} />
            <div className="relative z-0 text-gray-700 pt-32 pb-48 px-4">
                <h1 className="text-3xl font-light">Iniciar Sesión</h1>
            </div>
            <Footer />
        </>
    );
}
