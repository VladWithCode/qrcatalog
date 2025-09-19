import { Header } from '@/components/dashboard/header'
import { DashboardSidebar } from '@/components/dashboard/sidebar'
import { PageWrapper } from '@/components/pageWrapper'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated')({
    component: RouteComponent,
})

function RouteComponent() {
    return (
        <SidebarProvider>
            <DashboardSidebar />
            <SidebarInset>
                <PageWrapper className="grid-rows-[auto_1fr] bg-gray-200">
                    <Header />
                    <Outlet />
                </PageWrapper>
            </SidebarInset>
        </SidebarProvider>
    )
}
