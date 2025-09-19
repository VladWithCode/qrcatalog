import { useEffect, useState } from "react";

export function FloatingWhatsappBtn() {
    const [query, setQuery] = useState("");
    useEffect(() => {
        let usp = new URLSearchParams();
        usp.set(
            "text",
            "¡Hola!, me encantaría saber más sobre los servicios de QR Estrellas de la Limpieza.",
        );
        setQuery(usp.toString());
    }, []);

    return (
        <div id="wsp-float" className="fixed bottom-4 right-4 z-40">
            <a
                href={`https://wa.me/6184451290?${query}`}
                target="_blank"
                className="block rounded-full bg-wsp text-stone-50 p-3 shadow-sm shadow-stone-800/50"
            >
                <svg className="size-8 fill-current">
                    <use href="/wsp.svg#wsp"></use>
                </svg>
            </a>
        </div>
    );
}
