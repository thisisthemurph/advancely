import { cva, VariantProps } from "class-variance-authority";
import { cn } from "./lib/utils";

const loadingDotVariants = cva(
  "inline-block bg-current pointer-events-none aspect-square",
  {
    variants: {
      size: {
        xs: "w-4 h-4",
        sm: "w-5 h-5",
        md: "w-6 h-6",
        lg: "w-10 h-10",
      },
    },
    defaultVariants: {
      size: "md",
    },
  }
);

export interface LoadingDotProps
  extends VariantProps<typeof loadingDotVariants> {
  className?: string | undefined;
}

function LoadingDot({ className, size }: LoadingDotProps) {
  const imageUrl =
    "url(data:image/svg+xml;base64,Cjxzdmcgd2lkdGg9JzI0JyBoZWlnaHQ9JzI0JyB2aWV3Qm94PScwIDAgMjQgMjQnIHhtbG5zPSdodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2Zyc+PHN0eWxlPi5zcGlubmVyX3JYTlB7YW5pbWF0aW9uOnNwaW5uZXJfWWVCaiAuOHMgaW5maW5pdGV9QGtleWZyYW1lcyBzcGlubmVyX1llQmp7MCV7YW5pbWF0aW9uLXRpbWluZy1mdW5jdGlvbjpjdWJpYy1iZXppZXIoMC4zMywwLC42NiwuMzMpO2N5OjVweH00Ni44NzUle2N5OjIwcHg7cng6NHB4O3J5OjRweH01MCV7YW5pbWF0aW9uLXRpbWluZy1mdW5jdGlvbjpjdWJpYy1iZXppZXIoMC4zMywuNjYsLjY2LDEpO2N5OjIwLjVweDtyeDo0LjhweDtyeTozcHh9NTMuMTI1JXtyeDo0cHg7cnk6NHB4fTEwMCV7Y3k6NXB4fX08L3N0eWxlPjxlbGxpcHNlIGNsYXNzPSdzcGlubmVyX3JYTlAnIGN4PScxMicgY3k9JzUnIHJ4PSc0JyByeT0nNCcvPjwvc3ZnPg==)";

  return (
    <span
      style={{
        maskSize: "100%",
        WebkitMaskSize: "100%",
        maskPosition: "center",
        WebkitMaskPosition: "center",
        maskRepeat: "no-repeat",
        WebkitMaskRepeat: "no-repeat",
        maskImage: imageUrl,
        WebkitMaskImage: imageUrl,
      }}
      className={cn(loadingDotVariants({ size, className }))}
    ></span>
  );
}

export default LoadingDot;
