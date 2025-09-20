import type { TSectionParagraph } from "@/sections";

// Component to render paragraph content with list parsing
export function ParagraphContent({ content }: { content: string }) {
    const LIST_SEPARATOR = '<=>';

    if (content.includes(LIST_SEPARATOR)) {
        const items = content.split(LIST_SEPARATOR).filter(item => item.trim() !== '');
        return (
            <ul className="list-disc list-inside space-y-2">
                {items.map((item, index) => (
                    <li key={index} className="text-inherit">
                        {item.trim()}
                    </li>
                ))}
            </ul>
        );
    }

    // Regular paragraph content with line breaks
    return (
        <p className="text-inherit whitespace-pre-line translate-y-12 opacity-0" data-view-animate>
            {content}
        </p>
    );
}

// Component to render section paragraphs
export function SectionParagraphs({ paragraphs }: { paragraphs: TSectionParagraph[] }) {
    if (!paragraphs || paragraphs.length === 0) {
        return null;
    }

    return (
        <>
            {paragraphs.map((paragraph, index) => (
                <div key={paragraph.id || index} className="mb-4">
                    <ParagraphContent content={paragraph.content} />
                </div>
            ))}
        </>
    );
}

