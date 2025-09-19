import { Header } from "@/components/header";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useEffect, useRef } from "react";
import { useGSAP } from "@gsap/react";
import gsap from "gsap";
import { Footer } from "@/components/footer";
import { FloatingWhatsappBtn } from "@/components/wsp";
import { PageWrapper } from "@/components/pageWrapper";

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
        <PageWrapper ref={container}>
            <Header />
            <Outlet />
            <Footer />
            <FloatingWhatsappBtn />
        </PageWrapper>
    );
}
