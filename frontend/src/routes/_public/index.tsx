import { createFileRoute } from "@tanstack/react-router";
import { useSuspenseQuery } from "@tanstack/react-query";
import { publicSectionsQueryOptions, type TSection } from "@/sections";
import { AirbnbSection } from "@/components/dashboard/sections/airbnb-section";
import { SectionParagraphs } from "@/components/sections";

export const Route = createFileRoute("/_public/")({
    component: Index,
    loader: async ({ context }) => {
        return context.queryClient.ensureQueryData(publicSectionsQueryOptions);
    },
});

function Index() {
    const sections = useSuspenseQuery(publicSectionsQueryOptions).data;

    // Find sections by name for specific content areas
    const heroSection = sections.find(s => s.name === 'Sección Cabecera' || s.name === 'cabecera');
    const aboutSection = sections.find(s => s.name === 'Sección Nosotros' || s.name === 'nosotros');
    const servicesSection = sections.find(s => s.name === 'Sección Servicio' || s.name === 'servicio');
    const offerSection = sections.find(s => s.name === 'Sección Oferta' || s.name === 'oferta');
    const limpiezaAirbnbSection = sections.find(s => s.name === 'Sección Limpieza' || s.name === 'airbnb');
    const fumigationSection = sections.find(s => s.name === 'Sección Fumigación' || s.name === 'fumigacion');
    const laundrySection = sections.find(s => s.name === 'Sección Lavanderia' || s.name === 'lavanderia');
    const missionSection = sections.find(s => s.name === 'Sección Misión' || s.name === 'mision');
    const productsSection = sections.find(s => s.name === 'Sección Productos' || s.name === 'productos');
    const coverageSection = sections.find(s => s.name === 'Sección Cobertura' || s.name === 'cobertura');
    const visionSection = sections.find(s => s.name === 'Sección Visión' || s.name === 'vision');
    return (
        <>
            <section
                className="relative h-[95vh] z-0 overflow-hidden"
                id="inicio"
                data-section="inicio"
            >
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
                        {heroSection?.title || "Estrellas de la limpieza"}
                    </h1>
                    <div
                        className="text-current/80 md:text-2xl translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        <SectionParagraphs paragraphs={heroSection?.paragraphs || []} />
                    </div>

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
                        {aboutSection?.title || "¿Quiénes Somos?"}
                    </h2>
                    <div
                        className="translate-y-20 opacity-0"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        <SectionParagraphs paragraphs={aboutSection?.paragraphs || []} />
                    </div>
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
                        {servicesSection?.title || "Nuestro Servicio"}
                    </h2>
                    <div className="space-y-4">
                        <SectionParagraphs paragraphs={servicesSection?.paragraphs || []} />
                    </div>
                </div>
            </section>

            <AirbnbSection section={limpiezaAirbnbSection as TSection} />

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
                        {offerSection?.title || "Nuestra Oferta"}
                    </h2>
                    <div className="space-y-4">
                        <SectionParagraphs paragraphs={offerSection?.paragraphs || []} />
                    </div>
                </div>
            </section>

            <section
                id="fumigacion"
                className="relative z-0 py-16 px-4 md:px-12"
                data-section="fumigacion"
            >
                <div className="relative flex flex-col z-10 gap-4 xl:flex-row">
                    <div className="basis-1/2 shrink grow space-y-4 my-auto xl:order-2">
                        <h2
                            className="text-4xl font-bold opacity-0 translate-y-20"
                            data-view-animate="fadeIn"
                            data-view-animate-pos=">"
                        >
                            {fumigationSection?.title || "Servicio de Fumigación"}
                        </h2>
                        <div
                            className="font-medium opacity-0 translate-y-20"
                            data-view-animate="fadeIn"
                            data-view-animate-pos="<+=0.2"
                        >
                            <SectionParagraphs paragraphs={fumigationSection?.paragraphs || []} />
                        </div>
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

            <section
                id="lavanderia"
                className="relative py-16 px-4 z-0"
                data-section="lavanderia"
            >
                <div className="relative text-end flex flex-col z-10 gap-4 xl:flex-row">
                    <div className="basis-1/2 shrink grow space-y-4 my-auto">
                        <h2
                            className="text-4xl font-bold opacity-0 translate-y-20"
                            data-view-animate="fadeIn"
                            data-view-animate-pos=">"
                        >
                            {laundrySection?.title || "Servicio de Lavanderia"}
                        </h2>
                        <div
                            className="font-medium opacity-0 translate-y-20"
                            data-view-animate="fadeIn"
                            data-view-animate-pos="<+=0.2"
                        >
                            <SectionParagraphs paragraphs={laundrySection?.paragraphs || []} />
                        </div>
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
                        {productsSection?.title || "Productos"}
                    </h2>
                    <SectionParagraphs paragraphs={productsSection?.paragraphs || []} />
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
                        {missionSection?.title || "Misión"}
                    </h2>
                    <div
                        className="font-medium opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                    >
                        <SectionParagraphs paragraphs={missionSection?.paragraphs || []} />
                    </div>
                </div>
            </section>

            <section
                id="cobertura"
                className="relative z-0 space-y-4 py-16 px-4 md:px-12"
                data-section="cobertura"
            >
                <div className="flex flex-col text-end gap-4">
                    <h2
                        className="text-4xl font-bold opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos=">"
                    >
                        {coverageSection?.title || "Cobertura"}
                    </h2>
                    <SectionParagraphs paragraphs={coverageSection?.paragraphs.slice(0, 1) || []} />
                    <div className="relative bg-stone-400 aspect-[4/3] rounded-sm z-10 -mx-4 mt-4">
                        <img
                            src="/mapa_mex.webp"
                            alt=""
                            className="w-full h-full object-cover"
                        />
                    </div>
                    <div
                        className="font-medium opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                        data-view-animate-pos="<+=0.2"
                    >
                        <SectionParagraphs paragraphs={coverageSection?.paragraphs.slice(1) || []} />
                    </div>
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
                        {visionSection?.title || "Visión"}
                    </h2>
                    <div
                        className="font-medium opacity-0 translate-y-20"
                        data-view-animate="fadeIn"
                    >
                        <SectionParagraphs paragraphs={visionSection?.paragraphs || []} />
                    </div>
                </div>
            </section>




        </>
    );
}
