import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_public/")({
    component: Index,
});

function Index() {
    return (
        <>
            <section className="relative h-[95vh] z-0 overflow-hidden" id="inicio" data-section="inicio">
                <div className="absolute inset-0 z-0">
                    <video
                        className="h-full max-w-full object-cover brightness-50 md:w-full"
                        src="/qnr_main.mp4"
                        width="1280"
                        height="720"
                        autoPlay
                        loop
                        muted
                    ></video>
                </div>
                <div className="relative z-10 text-stone-50 h-full flex flex-col gap-6 items-center justify-center px-4 md:gap-8">
                    <h1
                        className="text-5xl font-medium space-y-4 md:text-8xl translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos=">"
                    >
                        Estrellas de la limpieza
                    </h1>
                    <p
                        className="text-current/80 md:text-2xl translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        El mejor servicio de limpieza en toda la República Mexicana.
                    </p>

                    <a
                        href="#contacto"
                        className="w-full max-w-96 bg-secondary-dark text-center py-4 rounded-full font-extrabold uppercase mt-18 md:py-8 md:text-2xl opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        ¡Contáctanos ya!
                    </a>
                </div>
            </section>
            <section
                className="relative flex flex-col py-16 px-4 gap-4 md:flex-row md:px-12 lg:px-24 lg:py-24 z-0"
                id="nosotros"
                data-section="nosotros"
            >
                <div className="space-y-4 my-auto">
                    <h2
                        className="text-4xl font-bold opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos=">"
                    >
                        ¿Quiénes Somos?
                    </h2>
                    <p
                        className="translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        Somos una empresa orgullosamente mexicana con presencia en distintos
                        estados de la república, distinguida por el uso de equipo de última
                        generación y personal altamente capacitado, obteniendo así la más alta
                        calidad en los servicios ofertados a un precio justo, cuidando siempre
                        tu economía.
                    </p>
                </div>
                <div
                    className="w-full bg-stone-400 aspect-video rounded-sm opacity-0 overflow-hidden"
                    data-view-animate="fadeIn"
                    data-view-animate-pos=">"
                >
                    <video
                        src="/limpieza_qnr.mp4"
                        className="w-full h-full max-w-full object-cover"
                        autoPlay
                        loop
                        muted
                    ></video>
                </div>
            </section>

            <section className="relative z-10 pb-24" id="limpieza" data-section="limpieza">
                <div className="absolute inset-0 z-0">
                    <img
                        className="h-full w-full object-cover object-center brightness-50"
                        src="/service_4.webp"
                        alt=""
                    />
                </div>
                <div className="relative z-10 text-stone-50 pt-16 pb-32 -mb-24 px-4 space-y-4">
                    <h2
                        className="text-4xl font-bold opacity-0 translate-y-20 text-stone-50"
                        data-view-animate="fadeIn"
                        data-view-animate-pos=">"
                    >
                        Nuestro Servicio
                    </h2>
                    <p
                        className="font-medium translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        QR Estrellas de la Limpieza ofrece un servicio integral y profesional
                        con presencia en diversas ciudades, abarcando un amplio territorio.
                    </p>
                    <p
                        className="font-medium translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        Contamos con documentación completa y actualizada, opciones flexibles de
                        pago y facturación a medida para adaptarnos a tus necesidades.{" "}
                    </p>
                    <p
                        className="font-medium translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        Además, nuestra póliza de responsabilidad civil respalda la seguridad y
                        confianza de nuestros clientes. Utilizamos maquinaria de alta gama y
                        productos de calidad superior, elaborados cuidadosamente para garantizar
                        resultados que no encontrarás en ningún otro lugar.
                    </p>
                    <p
                        className="font-medium translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        Disponibles las 24 horas del día, nos aseguramos de que tu espacio
                        siempre esté impecable. Nuestro compromiso es brindar calidad,
                        eficiencia y confianza, destacándonos como referentes en el sector.
                    </p>
                </div>
            </section>

            <section className="relative z-10 text-stone-50 text-end -mt-24" id="airbnb" data-section="airbnb">
                <svg
                    className="fill-primary-dark translate-y-px"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 1440 320"
                >
                    <path d="M0,64L34.3,106.7C68.6,149,137,235,206,245.3C274.3,256,343,192,411,160C480,128,549,128,617,106.7C685.7,85,754,43,823,58.7C891.4,75,960,149,1029,170.7C1097.1,192,1166,160,1234,144C1302.9,128,1371,128,1406,128L1440,128L1440,320L1405.7,320C1371.4,320,1303,320,1234,320C1165.7,320,1097,320,1029,320C960,320,891,320,823,320C754.3,320,686,320,617,320C548.6,320,480,320,411,320C342.9,320,274,320,206,320C137.1,320,69,320,34,320L0,320Z"></path>
                </svg>
                <div className="flex flex-col bg-primary-dark px-4 gap-4 py-8 md:px-12 md:py-16 lg:px-24 lg:flex-row">
                    <div className="basis-1/2 shrink grow space-y-4 xl:order-2 my-auto">
                        <h2
                            className="text-4xl font-bold opacity-0 translate-y-20"
                            data-view-animate="fadeIn"
                            data-view-animate-pos=">"
                        >
                            Limpieza
                            <span className="block text-secondary-dark">Airbnb</span>
                        </h2>
                        <p
                            className="font-medium translate-y-20 opacity-0"
                            data-view-animate="fadeIn"
                            data-view-animate-pos="<+=0.2"
                        >
                            En QR Estrellas de la Limpieza, nos dedicamos a ofrecer un servicio
                            de limpieza y desinfección de la más alta calidad para tu hogar,
                            oficina o vehículo.
                        </p>
                    </div>
                    <div
                        className="basis-1/2 shrink grow bg-stone-400 aspect-video rounded-sm overflow-hidden opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos=">"
                    >
                        <video
                            src="/qnr_cleaning.mp4"
                            className="w-full h-full max-w-full object-cover"
                            autoPlay
                            loop
                            muted
                        ></video>
                    </div>
                </div>
                <div className="grid grid-cols-3 gap-0.5 font-semibold bg-primary-dark p-4 pb-16 overflow-hidden">
                    <div className="col-span-full self-center justify-center mt-4">
                        <h3 className="text-center text-xl font-semibold tracking-wide uppercase py-4 text-stone-50 w-full">
                            Servicio Básico Casa Habitación
                        </h3>
                    </div>
                    <div className="col-start-1 col-span-2 bg-secondary-dark p-4 text-stone-50 rounded uppercase mb-1">
                        <ul className="text-sm text-start">
                            <li className="">- 1 Habitación Estándar</li>
                            <li className="">- 1 Baño</li>
                            <li className="">- 1 Sala, Comedor, Cocina</li>
                            <li className="">- 1 Pasillo y Estancia</li>
                        </ul>
                    </div>
                    <div className="col-span-1 text-xl bg-secondary-dark p-4 text-stone-50 uppercase rounded mb-1 flex items-center justify-center">
                        $500.00
                    </div>

                    <div className="col-start-1 col-span-2 border-2 border-stone-50 p-4 text-stone-50 rounded uppercase mb-1">
                        <ul className="text-sm text-start">
                            <li className="">- 2 Habitación Estandar</li>
                            <li className="">- 2 Baño</li>
                            <li className="">- 1 Sala, Comedor, Cocina</li>
                            <li className="">- 1 Pasillo y Estancia</li>
                        </ul>
                    </div>
                    <div className="col-span-1 text-xl p-4 text-stone-50 border-2 border-stone-50 uppercase rounded mb-1 flex items-center justify-center">
                        $600.00
                    </div>

                    <div className="col-start-1 col-span-2 bg-secondary-dark p-4 text-stone-50 rounded uppercase mb-1">
                        <ul className="text-sm text-start">
                            <li className="">- 3 Habitación Estandar</li>
                            <li className="">- 3 Baño</li>
                            <li className="">- 1 Sala, Comedor, Cocina</li>
                            <li className="">- 1 Pasillo y Estancia</li>
                        </ul>
                    </div>
                    <div className="col-span-1 text-xl bg-secondary-dark p-4 text-stone-50 uppercase rounded mb-1 flex items-center justify-center">
                        $700.00
                    </div>

                    <div className="col-start-1 col-span-2 border-2 border-stone-50 p-4 text-stone-50 rounded uppercase mb-1">
                        <ul className="text-sm text-start">
                            <li className="">- 4 Habitación Estandar</li>
                            <li className="">- 4 Baño</li>
                            <li className="">- 1 Sala, Comedor, Cocina</li>
                            <li className="">- 1 Pasillo y Estancia</li>
                        </ul>
                    </div>
                    <div className="col-span-1 text-xl p-4 text-stone-50 border-2 border-stone-50 uppercase rounded mb-1 flex items-center justify-center">
                        $800.00
                    </div>

                    <div className="col-span-full self-center justify-center mt-4 space-y-2">
                        <button
                            id="openPopupBtn"
                            className="w-full bg-secondary-dark text-stone-50 font-bold tracking-wide uppercase py-4 rounded-full"
                        >
                            Ver Detalles del Servicio
                        </button>
                        <button
                            id="openDetailsBtn"
                            className="w-full border-2 border-secondary-dark text-stone-50 font-bold tracking-wide uppercase py-4 rounded-full"
                        >
                            Ver Servicios Extra
                        </button>
                    </div>
                </div>
                <svg
                    className="fill-primary-dark rotate-180 -translate-y-px"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 1440 320"
                >
                    <path d="M0,64L34.3,106.7C68.6,149,137,235,206,245.3C274.3,256,343,192,411,160C480,128,549,128,617,106.7C685.7,85,754,43,823,58.7C891.4,75,960,149,1029,170.7C1097.1,192,1166,160,1234,144C1302.9,128,1371,128,1406,128L1440,128L1440,320L1405.7,320C1371.4,320,1303,320,1234,320C1165.7,320,1097,320,1029,320C960,320,891,320,823,320C754.3,320,686,320,617,320C548.6,320,480,320,411,320C342.9,320,274,320,206,320C137.1,320,69,320,34,320L0,320Z"></path>
                </svg>
            </section>

            <section className="relative z-0 pt-32 pb-16 -mt-24 px-4">
                <div className="absolute inset-0 z-0">
                    <img
                        className="h-full w-full object-cover brightness-50"
                        src="/service_1.webp"
                        alt=""
                    />
                </div>
                <div className="relative z-10 text-stone-50 space-y-4">
                    <h2
                        className="text-4xl font-bold opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos=">"
                    >
                        Nuestra Oferta
                    </h2>
                    <p
                        className="font-medium translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        QR Estrellas de la Limpieza ofrece una amplia gama de servicios
                        diseñados para cubrir todas las necesidades de nuestros clientes.
                    </p>
                    <p
                        className="font-medium translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        Desde autolavados y limpieza de casa habitación, hasta un servicio
                        express profesional 24 horas para Airbnb, adaptándonos a cada
                        requerimiento.
                    </p>
                    <p
                        className="font-medium translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        Además, contamos con lavanderías, venta de productos de limpieza de la
                        más alta calidad, elaborados a mano, y servicio especializado de lavado
                        de interiores de autobuses.
                    </p>
                    <p
                        className="font-medium translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        También ofrecemos fumigaciones y mantenimiento integral, todo con la
                        garantía de un trabajo impecable y atención personalizada. Con presencia
                        en varias ciudades y disponibilidad 24/7, estamos comprometidos con
                        brindar soluciones completas, eficientes y de calidad.
                    </p>
                </div>
            </section>

            <section id="fumigacion" className="relative z-0 py-16 px-4 md:px-12" data-section="fumigacion">
                <div className="relative flex flex-col z-10 gap-4 xl:flex-row">
                    <div className="basis-1/2 shrink grow space-y-4 my-auto xl:order-2">
                        <h2
                            className="text-4xl font-bold opacity-0 translate-y-20"
                            data-view-animate="fadeIn"
                            data-view-animate-pos=">"
                        >
                            Servicio de
                            <span className="text-secondary-dark">Fumigación</span>
                        </h2>
                        <p
                            className="font-medium opacity-0 translate-y-20"
                            data-view-animate="fadeIn"
                            data-view-animate-pos="<+=0.2"
                        >
                            Nuestra misión es erradicar cualquier tipo de plaga en los
                            domicilios, talleres y unidades de transporte para con esto asegurar
                            la seguridad y confianza de los pasajeros y empleados de los
                            talleres y oficinas donde se brindará el servicio.
                        </p>
                    </div>
                    <div
                        className="basis-1/2 shrink grow bg-stone-400 aspect-video rounded-sm overflow-hidden opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos=">"
                    >
                        <video
                            src="/fumi_vid.mov"
                            className="w-full h-full max-w-full object-cover"
                            autoPlay
                            loop
                            muted
                        ></video>
                    </div>
                </div>
                <div className="relative grid grid-cols-3 gap-0.5 text-center font-semibold py-4 overflow-hidden mt-4">
                    <div className="h-full flex items-center justify-center text-sm font-semibold bg-secondary-dark p-4 text-stone-50 rounded uppercase mb-1">
                        Concepto
                    </div>
                    <div className="h-full text-sm font-semibold bg-secondary-dark p-4 text-stone-50 uppercase rounded mb-1">
                        Aplicación Preventiva
                    </div>
                    <div className="h-full text-sm font-semibold p-4 text-stone-50 bg-secondary-dark rounded uppercase">
                        Aplicación Corre
                    </div>

                    <div className="text-sm font-bold p-4 rounded self-center">
                        Casa menor a 25m2
                    </div>
                    <div className="text-sm font-bold p-4 rounded self-center">$250.00</div>
                    <div className="text-sm font-bold p-4 rounded self-center">$350.00</div>

                    <div className="text-sm font-bold p-4 rounded self-center">
                        Comercio menor a 25m2
                    </div>
                    <div className="text-sm font-bold p-4 rounded self-center">$350.00</div>
                    <div className="text-sm font-bold p-4 rounded self-center">$450.00</div>
                </div>
            </section>

            <section id="lavanderia" className="relative py-16 px-4 z-0" data-section="lavanderia">
                <div className="relative text-end flex flex-col z-10 gap-4 xl:flex-row">
                    <div className="basis-1/2 shrink grow space-y-4 my-auto">
                        <h2
                            className="text-4xl font-bold opacity-0 translate-y-20"
                            data-view-animate="fadeIn"
                            data-view-animate-pos=">"
                        >
                            Servicio de
                            <span className="block text-primary-dark">Lavanderia</span>
                        </h2>
                        <p
                            className="font-medium opacity-0 translate-y-20"
                            data-view-animate="fadeIn"
                            data-view-animate-pos="<+=0.2"
                        >
                            Tu ropa limpia, fresca y lista a tiempo. Servicio de lavado, secado
                            y planchado con cuidado profesional. Entregamos con calidad y
                            puntualidad.
                        </p>
                    </div>
                    <div
                        className="basis-1/2 shrink grow bg-stone-400 aspect-video rounded-sm overflow-hidden opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        <video
                            src="/qnr_laundry.mp4"
                            className="w-full h-full max-w-full object-cover"
                            autoPlay
                            loop
                            muted
                        ></video>
                    </div>
                </div>
            </section>

            <section id="productos" className="relative flex flex-col gap-4 px-4 py-16 z-0">
                <div className="text-center space-y-4">
                    <h2
                        className="text-4xl font-bold opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos=">"
                    >
                        Nuestros Productos
                    </h2>
                    <p
                        className="font-medium opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        Productos profesionales personalizados para cada tipo de trabajo
                    </p>
                </div>
                <div className="grid grid-cols-4 grid-rows-3 w-full aspect-[4/3] gap-2 select-none">
                    <button
                        data-product-image-btn="/products_1.jpg"
                        className="col-start-1 row-start-1 col-end-3 row-end-3"
                    >
                        <img
                            data-view-animate="fadeIn"
                            data-view-animate-pos="<+=0.2"
                            data-product-image="/products_1.jpg"
                            src="/products_1.jpg"
                            alt=""
                            className="h-full aspect-square rounded-sm object-cover object-center brightness-60 hover:brightness-100 active:brightness-100 hover:rotate-5 active:rotate-5 transition-[rotate,filter] hover:cursor-pointer opacity-0 translate-y-20"
                        />
                    </button>
                    <button
                        data-product-image-btn="/products_2.jpg"
                        className="col-start-3 row-start-1 col-end-5 row-end-2"
                    >
                        <img
                            data-view-animate="fadeIn"
                            data-view-animate-pos="<+=0.2"
                            data-product-image="/products_2.jpg"
                            src="/products_2.jpg"
                            alt=""
                            className="w-full aspect-[2/1] rounded-sm object-cover object-center brightness-60 hover:brightness-100 active:brightness-100 hover:rotate-5 active:rotate-5 transition-[rotate,filter] hover:cursor-pointer opacity-0 translate-y-20"
                        />
                    </button>
                    <button
                        data-product-image-btn="/products_3.jpg"
                        className="col-start-3 row-start-2 col-end-5 row-end-4"
                    >
                        <img
                            data-view-animate="fadeIn"
                            data-view-animate-pos="<+=0.2"
                            data-product-image="/products_3.jpg"
                            src="/products_3.jpg"
                            alt=""
                            className="h-full aspect-square rounded-sm object-cover object-center brightness-60 hover:brightness-100 active:brightness-100 hover:rotate-5 active:rotate-5 transition-[rotate,filter] hover:cursor-pointer opacity-0 translate-y-20"
                        />
                    </button>
                    <button
                        data-product-image-btn="/products_4.jpg"
                        className="col-start-1 row-start-3 col-end-3 row-end-3"
                    >
                        <img
                            data-view-animate="fadeIn"
                            data-view-animate-pos="<+=0.2"
                            data-product-image="/products_4.jpg"
                            src="/products_4.jpg"
                            alt=""
                            className="w-full aspect-[2/1] rounded-sm object-cover object-center brightness-60 hover:brightness-100 active:brightness-100 hover:rotate-5 active:rotate-5 transition-[rotate,filter] hover:cursor-pointer opacity-0 translate-y-20"
                        />
                    </button>
                </div>
            </section>

            <div
                id="productsPopup"
                className="fixed inset-0 bg-black/75 flex items-center justify-center z-50"
                role="dialog"
                aria-modal="true"
                aria-labelledby="productsTitle"
                style={{ visibility: "hidden" }}
            >
                <button
                    id="closeProducts"
                    className="absolute top-4 right-4 text-white hover:text-gray-300 z-10 p-2 rounded-full hover:bg-black/50 transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-white"
                    aria-label="Close popup"
                >
                    <svg
                        className="w-8 h-8"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="M6 18L18 6M6 6l12 12"
                        ></path>
                    </svg>
                </button>

                <div
                    id="productsContainer"
                    className="relative w-full h-full flex items-center justify-center p-4 z-0"
                >
                    <img src="/products_1.jpg" alt="" id="productsDisplay" />
                </div>

                <h2 id="productsTitle" className="sr-only">
                    Galería de Productos
                </h2>
            </div>

            <section id="mision" className="relative z-0 space-y-4 py-16 px-4 md:px-12">
                <div className="absolute inset-0 z-0">
                    <img
                        className="h-full w-full object-cover brightness-50"
                        src="/service_5.webp"
                        alt=""
                    />
                </div>
                <div className="relative text-stone-50 flex flex-col gap-4 z-10">
                    <h2
                        className="text-4xl font-bold opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos=">"
                    >
                        Misión
                    </h2>
                    <p
                        className="font-medium opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                    >
                        Generar un ambiente de bienestar mediante la conservación de áreas
                        limpias y libres de plaga para asegurar la tranquilidad y salud de
                        nuestros clientes a un precio justo.
                    </p>
                </div>
            </section>

            <section id="cobertura" className="relative z-0 space-y-4 py-16 px-4 md:px-12" data-section="cobertura">
                <div className="flex flex-col text-end gap-4">
                    <h2
                        className="text-4xl font-bold opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos=">"
                    >
                        Cobertura
                    </h2>
                    <p
                        className="font-medium opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        Operamos en multiples estados de México. Descubre si estamos en tu
                        ciudad.
                    </p>
                    <div className="relative bg-stone-400 aspect-[4/3] rounded-sm z-10 -mx-4 mt-4">
                        <img
                            src="/mapa_mex.webp"
                            alt=""
                            className="w-full h-full object-cover"
                        />
                    </div>
                    <p
                        className="font-bold text-lg opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        Contamos con sucursales en:
                    </p>
                    <ul
                        className="font-medium opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<30%"
                    >
                        <li>Mazatlán, Sinaloa.</li>
                        <li>Durango, Durango.</li>
                        <li>Monterrey, Nuevo León.</li>
                        <li>Guadalajara, Jalisco.</li>
                        <li>México, CDMX.</li>
                        <li>Chihuahua, Chihuahua.</li>
                        <li>Tepic, Nayarit.</li>
                        <li>Tijuana, Baja California.</li>
                    </ul>
                </div>
            </section>

            <section
                id="vision"
                className="relative z-0 space-y-4 pt-20 pb-32 -mb-26 px-4 md:px-12"
            >
                <div className="absolute inset-0 z-0">
                    <img
                        className="h-full w-full object-cover object-center brightness-50"
                        src="/service_3.webp"
                        alt=""
                    />
                </div>
                <div className="relative text-stone-50 flex flex-col gap-4 z-10">
                    <h2
                        className="text-4xl font-bold opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos=">"
                    >
                        Visión
                    </h2>
                    <p
                        className="font-medium opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                    >
                        Posicionarnos como la empresa líder en el servicio profesional de
                        limpieza y control de plagas a nivel nacional superando las expectativas
                        de nuestros clientes más exigentes.
                    </p>
                </div>
            </section>
        </>
    );
}
