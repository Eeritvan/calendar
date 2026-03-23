import { useSortable } from "@dnd-kit/react/sortable";
import {RestrictToVerticalAxis} from '@dnd-kit/abstract/modifiers';
import {RestrictToWindow} from '@dnd-kit/dom/modifiers';
import type { Calendar } from "@/types";

interface Props {
  item: Calendar;
  index: number;
}

const CalendarItem = ({ item, index }: Props) => {
  const { ref } = useSortable({
    modifiers: [RestrictToVerticalAxis, RestrictToWindow],
    id: item.id,
    index
  });

  return (
    <li ref={ref} className="bg-amber-400 rounded-2xl h-10">
      {item.name} {item.id}
    </li>
  )
}

export default CalendarItem
