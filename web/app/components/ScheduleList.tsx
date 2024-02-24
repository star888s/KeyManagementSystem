import React from 'react';
import { format } from 'date-fns';
import { Schedule } from '../model';

interface ScheduleItemProps {
  schedule: Schedule;
  setSelectedSchedule: (schedule: Schedule) => void;
}

const ScheduleItem: React.FC<ScheduleItemProps> = ({ schedule, setSelectedSchedule }) => {
  // ISO形式の日付をJavaScriptのDateオブジェクトに変換
  const startTime = new Date(schedule.startTime);
  const endTime = new Date(schedule.endTime);

  // DateオブジェクトをHH:mm形式にフォーマット
  const formattedStartTime = format(startTime, 'HH:mm');
  const formattedEndTime = format(endTime, 'HH:mm');

  return (
    <div className='card' onClick={() => setSelectedSchedule(schedule)}>
      <h2 className='title'>{schedule.name}</h2>
      <div className='time'>{`${formattedStartTime} - ${formattedEndTime}`}</div>
      <div>{schedule.memo}</div>
    </div>
  );
};

export default ScheduleItem;
