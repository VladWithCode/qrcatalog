export function Footer() {
    return (
        <footer className="relative z-0">
            <svg
                className="fill-primary-dark translate-y-px"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 1440 320"
            >
                <path d="M0,128L48,128C96,128,192,128,288,128C384,128,480,128,576,149.3C672,171,768,213,864,197.3C960,181,1056,107,1152,74.7C1248,43,1344,53,1392,58.7L1440,64L1440,320L1392,320C1344,320,1248,320,1152,320C1056,320,960,320,864,320C768,320,672,320,576,320C480,320,384,320,288,320C192,320,96,320,48,320L0,320Z"></path>
            </svg>
            <div className="bg-primary-dark px-4 pt-8 pb-16 text-center text-stone-50 space-y-4">
                <div className="space-y-4 mb-10">
                    <img
                        src="/logo_1.webp"
                        alt="logo"
                        className="w-48 mx-auto"
                        style={{
                            filter: "drop-shadow(0px 10px 30px var(--color-secondary-dark))",
                        }}
                    />
                </div>
                <div className="space-y-2">
                    <h3 className="text-xl font-medium tracking-wide uppercase">Menu</h3>
                    <ul className="space-y-2 text-lg">
                        <li className="underline underline-offset-2">
                            <a className="block w-full" href="#inicio" data-nav-anchor>
                                Inicio
                            </a>
                        </li>
                        <li className="underline underline-offset-2">
                            <a className="block w-full" href="#nosotros" data-nav-anchor>
                                Acerca de Nosotros
                            </a>
                        </li>
                        <li className="underline underline-offset-2">
                            <a className="block w-full" href="#airbnb" data-nav-anchor>
                                Servicio de Limpieza
                            </a>
                        </li>
                        <li className="underline underline-offset-2">
                            <a className="block w-full" href="#lavanderia" data-nav-anchor>
                                Servicio de Lavanderia
                            </a>
                        </li>
                        <li className="underline underline-offset-2">
                            <a className="block w-full" href="#fumigacion" data-nav-anchor>
                                Servicio de Fumigación
                            </a>
                        </li>
                        <li className="underline underline-offset-2">
                            <a className="block w-full" href="#cobertura" data-nav-anchor>
                                Cobertura
                            </a>
                        </li>
                        <li className="underline underline-offset-2">
                            <a className="block w-full" href="#contacto">
                                Contáctanos
                            </a>
                        </li>
                    </ul>
                </div>
                <div className="space-y-4">
                    <h3 className="text-xl font-medium tracking-wide uppercase">Contacto</h3>
                    <ul className="flex items-center justify-center gap-6 text-lg">
                        <li className="underline underline-offset-2">
                            <a href="https://facebook.com/61577697835759/" target="_blank">
                                <svg className="h-8 w-8 fill-current">
                                    <use href="/fb.svg#fb"></use>
                                </svg>
                            </a>
                        </li>
                        <li className="underline underline-offset-2">
                            <a
                                href="https://www.instagram.com/qrestrellasdelalimpieza/"
                                target="_blank"
                            >
                                <svg className="h-9 w-9 fill-current">
                                    <use href="/ig.svg#ig"></use>
                                </svg>
                            </a>
                        </li>
                        <li className="underline underline-offset-2">
                            <a href="https://wa.me/6184451290" target="_blank">
                                <svg className="h-8 w-8 fill-current">
                                    <use href="/wsp.svg#wsp"></use>
                                </svg>
                            </a>
                        </li>
                    </ul>

                    <ul className="flex flex-col gap-2">
                        <li className="flex mx-auto gap-2">
                            <svg className="h-5 w-4 fill-current">
                                <use href="/fb.svg#fb"></use>
                            </svg>
                            <p className="font-medium">QR Las Estrellas de la Limpieza</p>
                        </li>
                        <li className="font-medium">669-106-2736</li>
                        <li className="font-medium">618-445-1290</li>
                        <li className="font-medium">
                            <p>Calle Salvador Nava 147</p>
                            <p>Zona Centro C.P. 34000</p>
                            <p>Durango, Durango. México</p>
                        </li>
                    </ul>

                    <button className="w-full bg-secondary-dark text-stone-50 font-bold tracking-wide uppercase py-4 rounded-full">
                        Contáctanos
                    </button>
                </div>
            </div>

            <div className="relative z-0">
                <div className="absolute inset-0 z-0">
                    <img
                        src="/limpieza_sm.jpg"
                        alt=""
                        className="h-full w-full object-cover brightness-50"
                    />
                </div>
                <div className="relative text-stone-50 z-10 pt-12 pb-24 px-4">
                    <p className="text-5xl font-bold">Estrellas de la limpieza</p>
                </div>
                <div className="relative z-10 pb-4 px-2">
                    <p className="text-xs text-center font-semibold text-stone-50/80">
                        Copyright © 2025 Limpieza Q&R
                    </p>
                </div>
            </div>
        </footer>
    );
}
