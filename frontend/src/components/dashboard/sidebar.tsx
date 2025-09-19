import { Link } from "@tanstack/react-router";
import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarGroupContent, SidebarGroupLabel, SidebarHeader, SidebarMenu, SidebarMenuItem, SidebarTrigger } from "../ui/sidebar";

export function DashboardSidebar() {
    return (
        <Sidebar variant="inset">
            <SidebarHeader className="flex-row items-center justify-between p-2.5">
                <SidebarTrigger />
                <h2 className="flex-1 font-semibold text-gray-700">Q&R Limpieza</h2>
            </SidebarHeader>
            <SidebarContent>
                <SidebarGroup>
                    <SidebarGroup>
                        <SidebarGroupLabel>Navegación</SidebarGroupLabel>
                        <SidebarGroupContent>
                            <SidebarMenu>
                                {sidebarItems.map((item) => (
                                    <SidebarMenuItem key={item.label}>
                                        <Link to={item.to}>{item.label}</Link>
                                    </SidebarMenuItem>
                                ))}
                            </SidebarMenu>
                        </SidebarGroupContent>
                    </SidebarGroup>
                </SidebarGroup>
            </SidebarContent>
            {/* <SidebarFooter> */}
            {/* </SidebarFooter> */}
        </Sidebar>
    );
}

const sidebarItems = [
    {
        label: "Panel",
        to: "/dashboard",
    },
    // Disabled for now
    // content management will appear in dashboards home page
    // {
    //     label: "Secciones",
    //     to: "/_authenticated/secciones",
    // },
];
