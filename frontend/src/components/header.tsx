import { Button } from "./ui/button";
import { Home, HomeIcon, Users } from "lucide-react";
import {
    NavigationMenu,
    NavigationMenuContent,
    NavigationMenuItem,
    NavigationMenuList,
    NavigationMenuTrigger,
} from "./ui/navigation-menu";
import { useIsMobile } from "@/hooks/use-mobile";
import { useEffect, useRef, useState } from "react";
import { useGSAP } from "@gsap/react";
import gsap from "gsap";
import { useLocation } from "@tanstack/react-router";
import { cn } from "@/lib/utils";
import { useInView } from "react-intersection-observer";

export function Header({ className, noAnimate, alwaysOpaque }: { noAnimate?: boolean, alwaysOpaque?: boolean } & React.HTMLAttributes<HTMLDivElement>) {
    const { inView, ref: headerThreshold } = useInView();
    const container = useRef<HTMLDivElement>(null);
    const isMobile = useIsMobile();
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    const { contextSafe } = useGSAP({ scope: container });
    const animateHeader = contextSafe((inView: boolean, isMenuOpen: boolean) => {
        if (noAnimate) {
            return;
        }

        const conditionalProps: gsap.TweenVars = {}
        if (inView && !isMenuOpen) {
            if (!alwaysOpaque) {
                conditionalProps.backgroundColor = "var(--c-bg-header-translucent)";
                conditionalProps.color = "var(--color-gray-50)";
            }

            gsap.to("header", {
                y: "0rem",
                opacity: 1,
                position: "absolute",
                ease: "power1.in",
                duration: 0.3,
                ...conditionalProps,
            });
        } else {
            if (!alwaysOpaque) {
                conditionalProps.backgroundColor = "var(--c-bg-header-opaque)";
                conditionalProps.color = "var(--color-gray-800)";
            }

            gsap.to("header", {
                y: "0rem",
                opacity: 1,
                position: "fixed",
                ease: "power1.in",
                duration: 0.3,
                ...conditionalProps,
            });
        }
    });

    useEffect(() => {
        animateHeader(inView, isMenuOpen);
    }, [inView, isMenuOpen]);

    return (
        <div className="relative z-30 w-screen" ref={container}>
            <header
                className={cn(
                    "absolute top-0 inset-x-0 h-14 text-stone-50 flex items-center z-20 px-2 md:px-12",
                    className,
                    noAnimate && "fixed",
                    alwaysOpaque && "bg-gray-50 text-stone-800 shadow",
                )}
                data-state={isMenuOpen ? "open" : "closed"}
            >
                <a
                    className={cn(
                        "text-stone-50 -translate-y-9 opacity-0 h-10",
                        noAnimate && "translate-y-0 opacity-100",
                    )}
                    data-no-animate={noAnimate}
                    data-always-opaque={alwaysOpaque}
                    data-view-animate="fadeIn"
                >
                    <img src="/logo_1.webp" alt="logo" className="h-full md:h-16" />
                </a>
                <Button
                    id="menu-toggler"
                    className={cn(
                        "relative ml-auto h-9 w-9 flex flex-col gap-1 opacity-0 data-[sticking=true]:text-gray-800",
                        noAnimate && "translate-y-0 opacity-100",
                        alwaysOpaque && "text-stone-800",
                    )}
                    data-sticking={!inView || isMenuOpen ? "true" : "false"}
                    data-view-animate="fadeIn"
                    data-no-animate={noAnimate}
                    data-always-opaque={alwaysOpaque}
                    onClick={() => setIsMenuOpen((isOpen) => !isOpen)}
                    variant="ghost"
                >
                    <div className="w-5 h-0.5 bg-current" data-toggler-line="1"></div>
                    <div className="w-5 h-0.5 bg-current" data-toggler-line="2"></div>
                    <div className="w-5 h-0.5 bg-current" data-toggler-line="3"></div>
                </Button>
            </header>
            <div className="absolute top-[65vh] inset-x-0" ref={headerThreshold}></div>
            {isMobile ? <MobileNavigationMenu isOpen={isMenuOpen} /> : <HeaderNavigationMenu />}
        </div>
    );
}

export function HeaderNavigationMenu() {
    return (
        <NavigationMenu className="hidden xl:block" viewport={false}>
            <NavigationMenuList>
                <NavigationMenuItem>
                    <NavigationMenuTrigger>Inicio</NavigationMenuTrigger>
                    <NavigationMenuContent>
                        <p>Inicio</p>
                    </NavigationMenuContent>
                </NavigationMenuItem>
            </NavigationMenuList>
        </NavigationMenu>
    );
}

function MobileNavigationMenu({ isOpen }: { isOpen: boolean }) {
    const location = useLocation();
    const [activeSection, setActiveSection] = useState<string | null>(null);
    const nav = useRef<HTMLDivElement>(null);
    const tl = useRef<gsap.core.Timeline>(null);
    const { contextSafe } = useGSAP(
        () => {
            tl.current = gsap.timeline({
                paused: true,
                defaults: { duration: 0.3, ease: "power1.in" },
            });
            tl.current.to('[data-menu-animate="hor"]', {
                x: "0rem",
                opacity: 1,
                stagger: 0.05,
            });
            tl.current.to(
                '[data-menu-animate="ver"]',
                {
                    y: "0rem",
                    opacity: 1,
                },
                "-=0.1",
            );
        },
        { scope: nav },
    );
    const animateIn = contextSafe((isOpen: boolean) => {
        if (!tl.current) {
            return;
        }

        if (isOpen) {
            tl.current.play();
        } else {
            tl.current.reverse();
        }
    });

    useEffect(() => {
        if (location.hash === "") {
            setActiveSection("inicio");
            return;
        }

        setActiveSection(location.hash);
    }, [location.hash]);

    useEffect(() => {
        animateIn(isOpen);
    }, [isOpen]);

    return (
        <nav
            id="menu"
            className="fixed top-0 left-0 right-0 bg-stone-50 text-stone-800 font-medium shadow shadow-stone-800/60 pb-8 md:px-12 z-40 transition-[translate,_opacity,_transform] -translate-y-12 data-[state=open]:translate-y-14 opacity-0 data-[state=open]:opacity-100 pointer-events-none data-[state=open]:pointer-events-auto overflow-hidden"
            data-state={isOpen ? "open" : "closed"}
            ref={nav}
        >
            <ul className="space-y-2 [&_a]:px-4 [&_a]:py-2 pt-4">
                {mainNavigationItems.map((item) => (
                    <li
                        key={item.label}
                        data-menu-animate="hor"
                        className={cn(
                            "hover:cursor-pointer translate-x-full opacity-0",
                            "#" + activeSection === item.to
                                ? "bg-primary-dark text-gray-50"
                                : "hover:bg-gray-800/20 active:bg-gray-800/20",
                        )}
                    >
                        <a
                            className="block w-full"
                            href={item.to}
                            data-nav-id={item.label}
                            data-nav-anchor
                        >
                            {item.label}
                        </a>
                    </li>
                ))}
                <li data-menu-animate="ver" className="mt-4 translate-y-full opacity-0 px-2">
                    <Button variant="ghost" asChild>
                        <a
                            className="block text-center font-bold w-full bg-warning text-stone-50 py-2 rounded-full uppercase"
                            href="#contacto"
                        >
                            Contáctanos
                        </a>
                    </Button>
                </li>
            </ul>
        </nav>
    );
}

const mainNavigationItems = [
    {
        label: "Inicio",
        to: "#inicio",
        icon: Home,
    },
    {
        label: "Acerca de Nosotros",
        to: "#nosotros",
        icon: Users,
    },
    {
        label: "Servicio de Limpieza",
        to: "#limpieza",
        icon: HomeIcon,
    },
    {
        label: "Servicio de Fumigación",
        to: "#fumigacion",
        icon: HomeIcon,
    },
    {
        label: "Servicio de Lavanderia",
        to: "#lavanderia",
        icon: HomeIcon,
    },
    {
        label: "Cobertura",
        to: "#cobertura",
        icon: HomeIcon,
    },
];
