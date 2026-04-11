type Stat = {
    label: string;
    value: string;
    subtext: string;
    icon: React.ComponentType;
    change: string;
    trend: "up" | "down";
    color: string;
    bgColor: string;
}

export type { Stat }