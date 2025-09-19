import { Footer } from "@/components/footer";
import { Header } from "@/components/header";
import { Button } from "@/components/ui/button";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { createFileRoute, redirect } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { PageWrapper } from "@/components/pageWrapper";
import { checkAuthOptions, signInOptions } from "@/auth";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";

export const Route = createFileRoute("/iniciar-sesion")({
    component: RouteComponent,
    staticData: {
        withOpaqueHeader: true,
    },
    beforeLoad: async ({ context }) => {
        const authData = await context.queryClient.fetchQuery(checkAuthOptions);
        if (authData && authData.isAuthenticated) {
            throw redirect({ to: '/dashboard' });
        }
    },
});

const formSchema = z.object({
    username: z.string().min(2, {
        message: "El nombre de usuario debe tener al menos 2 caracteres",
    }),
    password: z.string().min(4, {
        message: "La contraseña debe tener al menos 8 caracteres",
    }),
});

type LoginFormData = z.infer<typeof formSchema>;

function RouteComponent() {
    const [submitError, setSubmitError] = useState<string | null>(null);
    const form = useForm<LoginFormData>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            username: "",
            password: "",
        },
    });
    const signInMut = useMutation(signInOptions);
    const onSubmit = async (data: z.infer<typeof formSchema>) => {
        const response = await signInMut.mutateAsync(data);
        if (!response?.isAuthenticated) {
            setSubmitError(response?.message || "Usuario o contraseña incorrectos");
        }
    };

    return (
        <PageWrapper>
            <Header noAnimate={true} alwaysOpaque={true} />
            <div className="relative z-0 text-gray-700 py-32 px-4">
                <Card>
                    <CardHeader>
                        <CardTitle>
                            <h1 className="text-2xl">Iniciar Sesión</h1>
                        </CardTitle>
                        <CardDescription>
                            <p>Ingresa tu usuario y contraseña para continuar.</p>
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <Form {...form}>
                            <form
                                className="flex flex-col gap-4"
                                onSubmit={form.handleSubmit(onSubmit)}
                            >
                                <FormField
                                    control={form.control}
                                    name="username"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Usuario</FormLabel>
                                            <FormControl>
                                                <Input {...field} />
                                            </FormControl>
                                            <FormMessage className="text-xs" />
                                        </FormItem>
                                    )}
                                />
                                <FormField
                                    control={form.control}
                                    name="password"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Contraseña</FormLabel>
                                            <FormControl>
                                                <Input type="password" {...field} />
                                            </FormControl>
                                            <FormMessage className="text-xs" />
                                        </FormItem>
                                    )}
                                />
                                {submitError && (
                                    <FormMessage className="text-xs">
                                        {submitError}
                                    </FormMessage>
                                )}
                                <Button type="submit" className="bg-primary-dark">
                                    Iniciar Sesión
                                </Button>
                            </form>
                        </Form>
                    </CardContent>
                </Card>
            </div>
            <Footer />
        </PageWrapper>
    );
}
