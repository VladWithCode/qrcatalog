import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardAction, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import type { TSection, TSectionParagraph } from '@/sections';
import { Edit, Loader2 } from 'lucide-react';
import { SectionEditForm } from './sectionEditForm';
import { useGSAP } from '@gsap/react';
import gsap from 'gsap';
import { useRef } from 'react';

const LIST_SEPARATOR = '<=>';
const BRIEF_LENGTH = 240;

export function SectionCard({ section }: { section: TSection }) {
    const cardRef = useRef<HTMLDivElement>(null);
    const [dialogOpen, setDialogOpen] = useState(false);
    const [isEditing, setIsEditing] = useState(false);
    const [isUpdating, setIsUpdating] = useState(false);

    // GSAP animation for card entrance
    useGSAP(() => {
        if (cardRef.current) {
            gsap.from(cardRef.current, {
                opacity: 0,
                y: 20,
                duration: 0.5,
                ease: "power2.out"
            });
        }
    }, []);

    const handleEdit = () => {
        if (cardRef.current) {
            gsap.to(cardRef.current, {
                scale: 0.95,
                opacity: 0,
                duration: 0.2,
                ease: "power2.in",
                onComplete: () => setIsEditing(true)
            });
        } else {
            setIsEditing(true);
        }
    };

    const handleCancelEdit = () => {
        setIsEditing(false);
    };

    const handleSave = (formData: any) => {
        setIsUpdating(true);
        // Simulate a brief loading state for better UX
        setTimeout(() => {
            setIsEditing(false);
            setIsUpdating(false);
        }, 500);
    };

    if (isEditing) {
        return (
            <SectionEditForm
                section={section}
                onSave={handleSave}
                onCancel={handleCancelEdit}
            />
        );
    }

    if (isUpdating) {
        return <SectionCardSkeleton />;
    }

    return (
        <>
            <Card
                ref={cardRef}
                key={section.id}
                className="transition-all duration-200 hover:shadow-lg hover:scale-[1.02] cursor-pointer"
            >
                <CardHeader className="gap-6">
                    <CardTitle className="flex items-center justify-between gap-2">
                        <div className="flex flex-col items-start justify-between">
                            <span className="text-xl line-clamp-2">{section.title}</span>
                            <span className="text-sm font-medium text-muted-foreground">{section.name}</span>
                        </div>
                    </CardTitle>
                    <CardAction className="flex flex-col gap-2">
                        <Button variant="secondary" size="sm" onClick={handleEdit}>
                            <Edit className="mr-2 h-4 w-4" />
                            Editar
                        </Button>
                        {/* <Button variant="destructive" size="sm"> */}
                        {/*     <Trash className="mr-2 h-4 w-4" /> */}
                        {/*     Eliminar */}
                        {/* </Button> */}
                    </CardAction>
                </CardHeader>
                <CardContent>
                    <div className="">
                        {renderSectionBrief(section.paragraphs, () => setDialogOpen(true))}
                    </div>
                </CardContent>
            </Card>

            <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
                <DialogContent className="w-[95%] max-w-4xl max-h-[80vh] overflow-y-auto">
                    <DialogHeader>
                        <DialogTitle className="text-2xl">{section.title}</DialogTitle>
                        <DialogDescription className="text-base">
                            {section.name}
                        </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-6 mt-6">
                        {renderFullSectionContent(section.paragraphs)}
                    </div>
                </DialogContent>
            </Dialog>
        </>
    )
}



// Render brief version of section paragraphs with "ver más" button
function renderSectionBrief(paragraphs: TSectionParagraph[], onShowMore?: () => void): React.ReactNode {
    if (!paragraphs || paragraphs.length === 0) {
        return <p className="text-sm text-muted-foreground">Sin contenido</p>;
    }

    const briefContent: React.ReactNode[] = [];
    let remainingLength = BRIEF_LENGTH;
    let hasMoreContent = false;

    for (const paragraph of paragraphs) {
        if (remainingLength <= 0) {
            hasMoreContent = true;
            break;
        }

        const content = paragraph.content;
        if (!content) continue;

        // If content contains list separator, parse it as list
        if (content.includes(LIST_SEPARATOR)) {
            const items = content.split(LIST_SEPARATOR).filter(item => item.trim() !== '');
            const itemsToShow = items.slice(0, Math.ceil(remainingLength / 50)); // Estimate 50 chars per item

            briefContent.push(
                <ul key={`para-${paragraph.id}`} className="list-disc list-inside space-y-1 mb-2">
                    {itemsToShow.map((item, index) => (
                        <li key={index} className="text-sm">
                            {item.trim().substring(0, 60)}{item.trim().length > 60 ? '...' : ''}
                        </li>
                    ))}
                    {items.length > itemsToShow.length && (
                        <li className="text-sm text-muted-foreground">...</li>
                    )}
                </ul>
            );

            remainingLength -= itemsToShow.join('').length;
            if (items.length > itemsToShow.length) {
                hasMoreContent = true;
            }
        } else {
            // Regular paragraph
            const chunk = content.substring(0, remainingLength);
            const isTruncated = content.length > remainingLength;

            briefContent.push(
                <p key={`para-${paragraph.id}`} className="text-sm mb-2">
                    {chunk.split('\n').map((line, index) => (
                        <span key={index}>
                            {line}
                            {index < chunk.split('\n').length - 1 && <br />}
                        </span>
                    ))}
                    {isTruncated && '...'}
                </p>
            );

            remainingLength -= chunk.length;
            if (isTruncated) {
                hasMoreContent = true;
            }
        }
    }

    // Add "ver más" button if there's more content
    if (hasMoreContent && onShowMore) {
        briefContent.push(
            <Button
                key="show-more"
                variant="link"
                size="sm"
                className="p-0 h-auto text-sm text-primary hover:text-primary/80"
                onClick={onShowMore}
            >
                Ver más
            </Button>
        );
    }

    return briefContent;
}

// Render full section content for the dialog
function renderFullSectionContent(paragraphs: TSectionParagraph[]): React.ReactNode {
    if (!paragraphs || paragraphs.length === 0) {
        return <p className="text-muted-foreground">Sin contenido</p>;
    }

    return paragraphs.map((paragraph) => {
        const content = paragraph.content;
        if (!content) return null;

        // If content contains list separator, parse it as list
        if (content.includes(LIST_SEPARATOR)) {
            const items = content.split(LIST_SEPARATOR).filter(item => item.trim() !== '');
            return (
                <div key={`full-para-${paragraph.id}`} className="mb-6">
                    <ul className="list-disc list-inside space-y-2 text-base leading-relaxed">
                        {items.map((item, index) => (
                            <li key={index} className="text-foreground">
                                {item.trim()}
                            </li>
                        ))}
                    </ul>
                </div>
            );
        } else {
            // Regular paragraph
            return (
                <div key={`full-para-${paragraph.id}`} className="mb-6">
                    <p className="text-base leading-relaxed text-foreground whitespace-pre-line">
                        {content}
                    </p>
                </div>
            );
        }
    });
}

// Loading skeleton component
function SectionCardSkeleton() {
    const skeletonRef = useRef<HTMLDivElement>(null);

    useGSAP(() => {
        if (skeletonRef.current) {
            gsap.from(skeletonRef.current, {
                opacity: 0,
                scale: 0.95,
                duration: 0.3,
                ease: "power2.out"
            });
        }
    }, []);

    return (
        <Card ref={skeletonRef} className="w-full">
            <CardHeader className="gap-6">
                <CardTitle className="flex items-center justify-between gap-2">
                    <div className="flex flex-col items-start justify-between gap-2">
                        <div className="h-6 bg-muted animate-pulse rounded w-48"></div>
                        <div className="h-4 bg-muted animate-pulse rounded w-32"></div>
                    </div>
                </CardTitle>
                <div className="flex flex-col gap-2">
                    <div className="h-9 bg-muted animate-pulse rounded w-20"></div>
                    <div className="h-9 bg-muted animate-pulse rounded w-20"></div>
                </div>
            </CardHeader>
            <CardContent>
                <div className="space-y-3">
                    <div className="h-4 bg-muted animate-pulse rounded w-full"></div>
                    <div className="h-4 bg-muted animate-pulse rounded w-3/4"></div>
                    <div className="h-4 bg-muted animate-pulse rounded w-1/2"></div>
                    <div className="flex items-center gap-2 mt-4">
                        <Loader2 className="h-4 w-4 animate-spin text-primary" />
                        <span className="text-sm text-muted-foreground">Actualizando sección...</span>
                    </div>
                </div>
            </CardContent>
        </Card>
    );
}
