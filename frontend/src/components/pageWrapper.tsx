import { cn } from "@/lib/utils";
import { forwardRef, type PropsWithChildren } from "react";

export const PageWrapper = forwardRef(
    (
        props: PropsWithChildren & React.HTMLAttributes<HTMLDivElement>,
        ref: React.Ref<HTMLDivElement>,
    ) => {
        return (
            <div
                className={cn(
                    "relative z-0 grid grid-rows-[auto_1fr_auto] h-screen w-screen overflow-x-hidden scroll-smooth",
                    props.className,
                )}
                {...props}
                ref={ref}
            >
                {props.children}
            </div>
        );
    },
);
