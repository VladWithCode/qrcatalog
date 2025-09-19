import { Header } from "@/components/header";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useEffect, useRef } from "react";
import { useGSAP } from "@gsap/react";
import gsap from "gsap";
import { Footer } from "@/components/footer";
import { FloatingWhatsappBtn } from "@/components/wsp";

export const Route = createFileRoute("/_public")({
    component: RouteComponent,
});

function RouteComponent() {
    const obsv = useRef<IntersectionObserver>(null);
    const container = useRef<HTMLDivElement>(null);
    const { contextSafe } = useGSAP({ scope: container });
    const animateEnter = contextSafe((el: HTMLElement) => {
        gsap.to(el, {
            opacity: 1,
            y: "0rem",
            duration: 0.5,
            ease: "power3.in",
        });
    });

    useEffect(() => {
        obsv.current = new IntersectionObserver((ents, o) => {
            for (const ent of ents) {
                if (ent.isIntersecting) {
                    animateEnter(ent.target as HTMLElement);
                    o.unobserve(ent.target);
                }
            }
        });
        const els = document.querySelectorAll("[data-view-animate]");
        for (const el of els) {
            obsv.current.observe(el);
        }

        return () => {
            if (obsv.current) {
                obsv.current.disconnect();
            }
        };
    }, []);

    return (
        <div
            className="relative h-screen w-screen grid grid-rows-[auto_1fr_auto] z-0 overflow-x-hidden scroll-smooth"
            ref={container}
        >
            <Header />
            <Outlet />
            <Footer />
            <FloatingWhatsappBtn />
        </div>
    );
}
