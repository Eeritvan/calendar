import { useVirtualizer } from '@tanstack/react-virtual';
import { useEffect, useRef } from 'react';

const getWeekStart = (date: Date): Date => {
  const result = new Date(date);
  result.setDate(result.getDate() - result.getDay());
  result.setHours(0, 0, 0, 0);
  return result;
};

const generateWeekData = (weekIndex: number, referenceDate: Date): Array<string> => {
  const weekStart = getWeekStart(referenceDate);
  const targetWeekStart = new Date(weekStart);
  targetWeekStart.setDate(weekStart.getDate() + (weekIndex * 7));

  const weekData: Array<string> = [];
  for (let day = 0; day < 7; day++) {
    const currentDate = new Date(targetWeekStart);
    currentDate.setDate(targetWeekStart.getDate() + day);

    const monthDay = currentDate.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric'
    });
    const year = currentDate.getFullYear();
    weekData.push(`${monthDay}\n${year}`);
  }

  return weekData;
};

const getWeekDates = (weekIndex: number, referenceDate: Date): { start: Date; end: Date } => {
  const weekStart = getWeekStart(referenceDate);
  const targetWeekStart = new Date(weekStart);
  targetWeekStart.setDate(weekStart.getDate() + (weekIndex * 7));

  const targetWeekEnd = new Date(targetWeekStart);
  targetWeekEnd.setDate(targetWeekStart.getDate() + 6);

  return { start: targetWeekStart, end: targetWeekEnd };
};

const COLUMN_WIDTH = 150;
const ROW_HEIGHT = 120;
const TOTAL_WEEKS = 2000;
const CURRENT_WEEK_INDEX = Math.floor(TOTAL_WEEKS / 2);

const Calendar = () => {
  const parentRef = useRef(null);
  const referenceDate = useRef(new Date()).current;

  const rowVirtualizer = useVirtualizer({
    count: TOTAL_WEEKS,
    getScrollElement: () => parentRef.current,
    estimateSize: () => (ROW_HEIGHT),
    initialOffset: CURRENT_WEEK_INDEX * ROW_HEIGHT,
    overscan: 5,
  });

  const columnVirtualizer = useVirtualizer({
    horizontal: true,
    count: 7,
    getScrollElement: () => parentRef.current,
    estimateSize: () => COLUMN_WIDTH,
  });

  const rowItems = rowVirtualizer.getVirtualItems();
  const columnItems = columnVirtualizer.getVirtualItems();

  useEffect(() => {
    if (rowItems.length > 0) {
      const firstRowIndex = rowItems[0].index;
      const lastRowIndex = rowItems[rowItems.length - 1].index;

      const firstWeekDates = getWeekDates(firstRowIndex - 1 - CURRENT_WEEK_INDEX, referenceDate);
      const lastWeekDates = getWeekDates(lastRowIndex - 1 - CURRENT_WEEK_INDEX, referenceDate);

      console.log('First visible date:', firstWeekDates.start.toLocaleDateString('en-US', {
        weekday: 'short',
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      }));
      console.log('Last visible date:', lastWeekDates.end.toLocaleDateString('en-US', {
        weekday: 'short',
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      }));
    }
  }, [rowItems, referenceDate]);

  return (
    <div
      ref={parentRef}
      className="h-150 overflow-auto bg-white"
    >
      <div
        className="relative"
        style={{
          height: `${rowVirtualizer.getTotalSize()}px`,
          width: `${columnVirtualizer.getTotalSize()}px`,
        }}
      >
        {rowItems.map((row) => {
          const weekData = generateWeekData(row.index - 1 - CURRENT_WEEK_INDEX, referenceDate);

          return (
            <div
              key={row.key}
              data-index={row.index}
              ref={rowVirtualizer.measureElement}
              className="absolute w-full grid grid-cols-7"
              style={{
                transform: `translateY(${row.start}px)`,
              }}
            >
              {columnItems.map((column) => {
                const cellContent = weekData[column.index];

                return (
                  <div
                    key={column.key}
                    className={`overflow-hidden whitespace-pre-wrap border hover:bg-red-200`}
                    style={{ minHeight: ROW_HEIGHT }}
                  >
                    {cellContent}
                  </div>
                );
              })}
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default Calendar;
