import { useState } from 'react';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';

export function AirbnbSection() {
    const [serviceDetailsOpen, setServiceDetailsOpen] = useState(false);
    const [extraServicesOpen, setExtraServicesOpen] = useState(false);

    return (
        <>
            <section
                id="airbnb"
                className="relative z-10 text-stone-50 text-end -mt-24"
                data-section="airbnb"
            >
                <svg
                    className="fill-primary-dark translate-y-px"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 1440 320"
                >
                    <path d="M0,64L34.3,106.7C68.6,149,137,235,206,245.3C274.3,256,343,192,411,160C480,128,549,128,617,106.7C685.7,85,754,43,823,58.7C891.4,75,960,149,1029,170.7C1097.1,192,1166,160,1234,144C1302.9,128,1371,128,1406,128L1440,128L1440,320L1405.7,320C1371.4,320,1303,320,1234,320C1165.7,320,1097,320,1029,320C960,320,891,320,823,320C754.3,320,686,320,617,320C548.6,320,480,320,411,320C342.9,320,274,320,206,320C137.1,320,69,320,34,320L0,320Z"></path>
                </svg>
                <div className="flex flex-col bg-primary-dark px-4 gap-4 py-8 md:px-12 md:py-16 lg:px-24">
                    <div className="basis-1/2 shrink grow space-y-4 my-auto xl:order-2">
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
                            className="w-full bg-secondary-dark text-stone-50 font-bold tracking-wide uppercase py-4 rounded-full hover:bg-secondary-dark/90 transition-colors"
                            onClick={() => setServiceDetailsOpen(true)}
                        >
                            Ver Detalles del Servicio
                        </button>
                        <button
                            className="w-full border-2 border-secondary-dark text-stone-50 font-bold tracking-wide uppercase py-4 rounded-full hover:bg-secondary-dark/10 transition-colors"
                            onClick={() => setExtraServicesOpen(true)}
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

            {/* Service Details Dialog */}
            <Dialog open={serviceDetailsOpen} onOpenChange={setServiceDetailsOpen}>
                <DialogContent className="max-w-4xl max-h-[90vh] overflow-hidden">
                    <DialogHeader>
                        <DialogTitle className="text-2xl text-center">
                            Detalles del Servicio de Limpieza Airbnb
                        </DialogTitle>
                        <DialogDescription className="text-center">
                            Información completa sobre nuestros servicios de limpieza profesional
                        </DialogDescription>
                    </DialogHeader>
                    <div className="flex items-center justify-center p-4">
                        <img
                            src="/details_service.jpeg"
                            alt="Detalles del servicio de limpieza Airbnb"
                            className="max-w-full max-h-[70vh] object-contain rounded-lg shadow-2xl"
                        />
                    </div>
                </DialogContent>
            </Dialog>

            {/* Extra Services Dialog */}
            <Dialog open={extraServicesOpen} onOpenChange={setExtraServicesOpen}>
                <DialogContent className="max-w-4xl max-h-[90vh] overflow-hidden">
                    <DialogHeader>
                        <DialogTitle className="text-2xl text-center">
                            Servicios Extra Disponibles
                        </DialogTitle>
                        <DialogDescription className="text-center">
                            Descubre nuestros servicios adicionales y opciones personalizadas
                        </DialogDescription>
                    </DialogHeader>
                    <div className="flex items-center justify-center p-4">
                        <img
                            src="/extra_service.jpeg"
                            alt="Servicios extra y opciones adicionales"
                            className="max-w-full max-h-[70vh] object-contain rounded-lg shadow-2xl"
                        />
                    </div>
                </DialogContent>
            </Dialog>
        </>
    );
}
