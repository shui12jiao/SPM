package db

const (
	// 预约状态('reserved', 'completed', 'canceled', 'violated')
	ReservationStatusReserved  = "reserved"  // 预约中
	ReservationStatusCompleted = "completed" // 已完成
	ReservationStatusCanceled  = "canceled"  // 已取消
	ReservationStatusViolated  = "violated"  // 违约
)
