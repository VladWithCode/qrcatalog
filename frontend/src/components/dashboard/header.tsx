import { SidebarTrigger } from "../ui/sidebar";

export function Header() {
    return (
        <div className="relative z-30 w-full px-0.5 py-1">
            <header className="flex items-center gap-6 justify-between p-2 rounded-lg shadow bg-gray-100">
                <SidebarTrigger />
                <h2 className="font-semibold text-gray-700 flex-1">
                    Panel de control
                </h2>
            </header>
        </div>
    )
}
