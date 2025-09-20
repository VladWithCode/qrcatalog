import { useForm, useFieldArray } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import type { TSection } from '@/sections';
import { updateSectionMutationOptions, sectionsQueryKeys } from '@/sections';
import { Save, X, Loader2, AlertCircle } from 'lucide-react';
import { useGSAP } from '@gsap/react';
import gsap from 'gsap';
import { useRef, useState } from 'react';

// Zod validation schema with Spanish messages
const sectionFormSchema = z.object({
    name: z.string()
        .min(1, { message: "El nombre es obligatorio" })
        .max(100, { message: "El nombre no puede tener más de 100 caracteres" })
        .regex(/^[\p{L}\p{N}_ ]*$/u, {
            message: "El nombre solo puede contener letras, números, espacios, y guiones bajos"
        }),
    title: z.string()
        .min(1, { message: "El título es obligatorio" })
        .max(200, { message: "El título no puede tener más de 200 caracteres" }),
    paragraphs: z.array(z.object({
        id: z.string().optional(),
        content: z.string()
            .min(1, { message: "El contenido del párrafo no puede estar vacío" })
            .max(2000, { message: "El párrafo no puede tener más de 2000 caracteres" })
    }))
        .min(1, { message: "Debe haber al menos un párrafo" })
        .max(10, { message: "No puede haber más de 10 párrafos" })
});

type SectionFormData = z.infer<typeof sectionFormSchema>;

interface SectionEditFormProps {
    section: TSection;
    onSave: (data: SectionFormData) => void;
    onCancel: () => void;
    isLoading?: boolean;
}

export function SectionEditForm({ section, onSave, onCancel }: SectionEditFormProps) {
    const formRef = useRef<HTMLDivElement>(null);
    const [submitError, setSubmitError] = useState<string | null>(null);

    const queryClient = useQueryClient();

    const form = useForm<SectionFormData>({
        resolver: zodResolver(sectionFormSchema),
        defaultValues: {
            name: section.name || '',
            title: section.title || '',
            paragraphs: section.paragraphs?.map(p => ({
                id: p.id,
                content: p.content || ''
            })) || [{ content: '' }]
        }
    });

    const { fields, append, remove } = useFieldArray({
        control: form.control,
        name: "paragraphs"
    });

    // GSAP animation for form entrance
    useGSAP(() => {
        if (formRef.current) {
            gsap.from(formRef.current, {
                opacity: 0,
                y: 20,
                duration: 0.5,
                ease: "power2.out"
            });
        }
    }, []);

    // Mutation for updating section
    const updateSectionMutation = useMutation({
        ...updateSectionMutationOptions(section.id),
        onSuccess: (data) => {
            // Invalidate and refetch sections
            queryClient.invalidateQueries({ queryKey: sectionsQueryKeys.all() });

            // Animate success and call onSave
            if (formRef.current) {
                gsap.to(formRef.current, {
                    scale: 1.02,
                    duration: 0.2,
                    yoyo: true,
                    repeat: 1,
                    ease: "power2.inOut",
                    onComplete: () => {
                        onSave(data.section);
                    }
                });
            } else {
                onSave(data.section);
            }
        },
        onError: (error: Error) => {
            setSubmitError(error.message || 'Ocurrió un error al guardar los cambios');
            // Animate error shake
            if (formRef.current) {
                gsap.to(formRef.current, {
                    x: -10,
                    duration: 0.1,
                    repeat: 5,
                    yoyo: true,
                    ease: "power2.inOut"
                });
            }
        }
    });

    const onSubmit = async (data: SectionFormData) => {
        setSubmitError(null);

        // Validate form before sending
        const isValid = await form.trigger();
        if (!isValid) {
            return;
        }

        // Transform data to match API expectations
        const updateData = {
            id: section.id,
            name: data.name,
            title: data.title,
            paragraphs: data.paragraphs.map((p, index) => ({
                id: p.id,
                section_id: section.id,
                order: index,
                content: p.content,
                content_as_list: false,
                created_at: '',
                updated_at: ''
            }))
        };

        updateSectionMutation.mutate(updateData as any);
    };

    const addParagraph = () => {
        if (fields.length < 10) {
            append({ content: '' });
            // Animate new field
            setTimeout(() => {
                const newFields = document.querySelectorAll('[data-field="paragraph"]');
                if (newFields.length > 0) {
                    gsap.from(newFields[newFields.length - 1], {
                        opacity: 0,
                        y: -10,
                        duration: 0.3,
                        ease: "power2.out"
                    });
                }
            }, 0);
        }
    };

    const removeParagraph = (index: number) => {
        if (fields.length > 1) {
            // Animate removal
            const fieldToRemove = document.querySelector(`[data-field="paragraph-${index}"]`);
            if (fieldToRemove) {
                gsap.to(fieldToRemove, {
                    opacity: 0,
                    y: -10,
                    duration: 0.2,
                    ease: "power2.in",
                    onComplete: () => remove(index)
                });
            } else {
                remove(index);
            }
        }
    };

    return (
        <Card ref={formRef} className="w-full">
            <CardHeader>
                <CardTitle className="flex items-center justify-between">
                    <span>Editar Sección</span>
                    <Button
                        variant="ghost"
                        size="sm"
                        onClick={onCancel}
                        disabled={updateSectionMutation.isPending}
                    >
                        <X className="h-4 w-4" />
                    </Button>
                </CardTitle>
            </CardHeader>
            <CardContent>
                {/* Error Display */}
                {submitError && (
                    <div className="mb-6 p-4 bg-destructive/10 border border-destructive/20 rounded-lg flex items-center gap-3">
                        <AlertCircle className="h-5 w-5 text-destructive flex-shrink-0" />
                        <p className="text-sm text-destructive">{submitError}</p>
                    </div>
                )}

                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                        {/* Section Name */}
                        <FormField
                            control={form.control}
                            name="name"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Nombre de la Sección</FormLabel>
                                    <FormControl>
                                        <Input
                                            placeholder="ej: servicios-principales"
                                            {...field}
                                            disabled={updateSectionMutation.isPending}
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        {/* Section Title */}
                        <FormField
                            control={form.control}
                            name="title"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Título de la Sección</FormLabel>
                                    <FormControl>
                                        <Input
                                            placeholder="ej: Nuestros Servicios Principales"
                                            {...field}
                                            disabled={updateSectionMutation.isPending}
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        {/* Paragraphs */}
                        <div className="space-y-4">
                            <div className="flex items-center justify-between">
                                <FormLabel className="text-base">Contenido de los Párrafos</FormLabel>
                                <Button
                                    type="button"
                                    variant="outline"
                                    size="sm"
                                    onClick={addParagraph}
                                    disabled={updateSectionMutation.isPending || fields.length >= 10}
                                >
                                    Agregar Párrafo
                                </Button>
                            </div>

                            {fields.map((field, index) => (
                                <FormField
                                    key={field.id}
                                    control={form.control}
                                    name={`paragraphs.${index}.content`}
                                    render={({ field: inputField }) => (
                                        <FormItem data-field={`paragraph-${index}`}>
                                            <div className="flex items-center justify-between">
                                                <FormLabel className="text-sm font-medium">
                                                    Párrafo {index + 1}
                                                </FormLabel>
                                                {fields.length > 1 && (
                                                    <Button
                                                        type="button"
                                                        variant="ghost"
                                                        size="sm"
                                                        onClick={() => removeParagraph(index)}
                                                        disabled={updateSectionMutation.isPending}
                                                        className="text-destructive hover:text-destructive"
                                                    >
                                                        <X className="h-4 w-4" />
                                                    </Button>
                                                )}
                                            </div>
                                            <FormControl>
                                                <Textarea
                                                    placeholder={`Contenido del párrafo ${index + 1}...`}
                                                    className="min-h-[100px] resize-y"
                                                    {...inputField}
                                                    disabled={updateSectionMutation.isPending}
                                                />
                                            </FormControl>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />
                            ))}
                        </div>

                        {/* Form Actions */}
                        <div className="flex justify-end gap-3 pt-4 border-t">
                            <Button
                                type="button"
                                variant="outline"
                                onClick={onCancel}
                                disabled={updateSectionMutation.isPending}
                            >
                                Cancelar
                            </Button>
                            <Button
                                type="submit"
                                disabled={updateSectionMutation.isPending}
                            >
                                {updateSectionMutation.isPending ? (
                                    <>
                                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                        Guardando...
                                    </>
                                ) : (
                                    <>
                                        <Save className="mr-2 h-4 w-4" />
                                        Guardar Cambios
                                    </>
                                )}
                            </Button>
                        </div>
                    </form>
                </Form>
            </CardContent>
        </Card>
    );
}
