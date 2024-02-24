'use client';
import React, { useEffect, useState } from 'react';
import 'react-datetime/css/react-datetime.css';
import ScheduleItem from './components/ScheduleList';
import fetchData from './api/fetchData';
import { Schedule, Info } from './model';
import SearchForm from './components/SearchForm';
import AddScheduleButton from './components/AddScheduleButton';
import ScheduleModal from './components/ShowModal';
import ScheduleDetailModal from './components/ScheduleDetailModal';
import fetchInfo from './api/fetchInfo';
import upsertSchedule from './api/upsertSchedule';
import deleteSchedule from './api/deleteSchedule';
import { info } from 'console';
import { Carousel } from 'react-responsive-carousel';
import 'react-responsive-carousel/lib/styles/carousel.min.css';

export default function Home() {
  // 取得
  const [schedules, setSchedules] = useState<Schedule[]>([]);

  const fetchAndSetSchedules = async () => {
    const schedules = await fetchData();
    setSchedules(schedules);
  };

  // 検索
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedSchedule, setSelectedSchedule] = useState<Schedule | null>(null);

  // 検索ワードに一致するスケジュールのみを表示する
  const filteredSchedules = schedules.filter(
    (schedule) => schedule.name.includes(searchTerm) || schedule.memo.includes(searchTerm),
  );

  // スケジュールを日付でソート
  const sortedSchedules = [...filteredSchedules].sort((a, b) => a.startTime.localeCompare(b.startTime));

  // 同じ日付のスケジュールをグループ化
  const schedulesByDate = sortedSchedules.reduce<Record<string, Schedule[]>>((groups, schedule) => {
    const date = schedule.startTime.split('T')[0];
    if (!groups[date]) {
      groups[date] = [];
    }
    groups[date].push(schedule);
    return groups;
  }, {});

  //登録
  // 登録用一覧取得
  const [infoList, setInfoList] = useState<Info[]>([]);

  const fetchAndSetInfo = async () => {
    const infoItems = await fetchInfo();
    setInfoList(infoItems);

    // infoItemsが空でない場合、newScheduleを更新します
    if (infoItems.length > 0) {
      const firstInfo = infoItems[0];
      setNewSchedule((prevSchedule) => ({ ...prevSchedule, id: firstInfo.id, name: firstInfo.name }));
    }
  };

  // 登録
  const [showModal, setShowModal] = useState(false);
  const [newSchedule, setNewSchedule] = useState({
    id: '',
    name: '',
    startTime: '',
    endTime: '',
    memo: '',
  });

  const handleIdChange = (e: { target: { value: any } }) => {
    console.log(infoList);
    const selectedInfo = infoList.find((info) => info.id === e.target.value);
    setNewSchedule({ ...newSchedule, id: e.target.value, name: selectedInfo?.name || '' });
  };

  const handleTimeChange = (e: { target: { name: any; value: any } }) => {
    setNewSchedule({ ...newSchedule, [e.target.name]: e.target.value + ':00+09:00' });
  };

  const handleInputChange = (e: { target: { name: any; value: any } }) => {
    setNewSchedule({ ...newSchedule, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: { preventDefault: () => void }) => {
    e.preventDefault();

    const response = await upsertSchedule(newSchedule);

    if (response == true) {
      fetchAndSetSchedules();
      setShowModal(false);
      alert('登録完了しました。');
    } else {
      alert('Failed to add schedule');
    }
  };

  // 削除
  const deleteScheduleHandler = async (schedule: Schedule) => {
    if (window.confirm('Are you sure you want to delete this schedule?')) {
      const response = await deleteSchedule(schedule);

      if (response == true) {
        fetchAndSetSchedules();
        setSelectedSchedule(null);
        alert('削除しました。');
      } else {
        alert('Failed to delete schedule');
      }
    }
  };

  //日付
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const options: Intl.DateTimeFormatOptions = { month: '2-digit', day: '2-digit', weekday: 'short' };
    return date.toLocaleDateString('ja-JP', options);
  };

  useEffect(() => {
    fetchAndSetSchedules();
    fetchAndSetInfo();
  }, []);

  return (
    <React.Fragment>
      <div className='navbar'>
        <div className='mb-8 space-between' style={{ padding: '10px' }}>
          <div style={{ display: 'flex', justifyContent: 'center', flex: 1 }}>
            <SearchForm searchTerm={searchTerm} setSearchTerm={setSearchTerm} />
          </div>
          <div>
            <AddScheduleButton setShowModal={setShowModal} />
          </div>
          <ScheduleModal
            showModal={showModal}
            setShowModal={setShowModal}
            newSchedule={newSchedule}
            handleIdChange={handleIdChange}
            handleTimeChange={handleTimeChange}
            handleInputChange={handleInputChange}
            handleSubmit={handleSubmit}
            infoList={infoList}
          />
        </div>
      </div>
      <div className='dashboard flex items-center justify-center'>
        <div style={{ display: 'flex', flexDirection: 'column', height: '80vh', width: '80vw' }}>
          <div style={{ overflowY: 'auto', flex: '1 1 auto' }}>
            <Carousel>
              {Object.entries(schedulesByDate).map(([date, schedules], index) => (
                <div key={index}>
                  <h2 style={{ fontSize: '36px' }}>{formatDate(date)}</h2>
                  {schedules.map((schedule, index) => (
                    <ScheduleItem key={index} schedule={schedule} setSelectedSchedule={setSelectedSchedule} />
                  ))}
                </div>
              ))}
            </Carousel>
          </div>
          <ScheduleDetailModal
            selectedSchedule={selectedSchedule}
            setSelectedSchedule={setSelectedSchedule}
            deleteSchedule={deleteScheduleHandler}
          />
        </div>
      </div>
    </React.Fragment>
  );
}
