import { SectionCard } from '@/components/dashboard/sections/sectionCard';
import { sectionsQueryOptions, type TSectionFilterParams } from '@/sections';
import { useSuspenseQuery } from '@tanstack/react-query';
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/dashboard')({
    component: RouteComponent,
    loader: async ({ context }) => {
        const defaultFilters: TSectionFilterParams = {
            page: 1,
            limit: 10,
        };
        return context.queryClient.ensureQueryData(sectionsQueryOptions(defaultFilters));
    },
})

function RouteComponent() {
    const { data: sections } = useSuspenseQuery(sectionsQueryOptions());
    return (
        <div className="space-y-6 px-2.5 py-6">
            <div className="px-3.5">
                <h1 className="text-3xl font-bold tracking-tight">Gestión de Contenido</h1>
                <p className="text-muted-foreground">
                    Administra el contenido de las diferentes secciones de tu sitio web.
                </p>
            </div>

            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {sections.sections.map((section) => (
                    <SectionCard key={section.id} section={section} />
                ))}
            </div>
        </div >
    );
}
