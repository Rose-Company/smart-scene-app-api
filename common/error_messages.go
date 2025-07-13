package common

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

const (
	TokenNotFound     = "TokenNotFound"
	TokenUnAuthorized = "TokenUnAuthorized"
	PerDenied         = "PerDenied"
	ResultFailed      = "ResultFailed"
	InternalError     = "InternalError"
)

const (
	DuplicatedDataErr = "Lỗi! Trùng lặp dữ liệu"
	DefaultError      = "Lỗi! Xảy ra lỗi không xác định, vui lòng liên hệ quản trị viên"
)

var DataIsNullErr = func(obj string) string {
	return fmt.Sprintf("%v cannot use nil", obj)
}

var DataIsExisted = func(obj string) string {
	return fmt.Sprintf("%v is existed", obj)
}

var DataIsSmallerZero = func(obj string) string {
	return fmt.Sprintf("%v is not smaller zero", obj)
}

var DataIsBeforeNow = func(obj string) string {
	return fmt.Sprintf("%v is not before now", obj)
}

var ErrorWrapper = func(prefix string, err error) error {
	return fmt.Errorf("%v: %v", prefix, err.Error())
}

var PgErrorTransform = func(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), "duplicate key value") {
		return fmt.Errorf(DuplicatedDataErr)
	}

	return err
}

var (
	ErrCodeNotAuthorized                            = errors.New("not_authorized")
	ErrCartCodeNotFound                             = errors.New("cart_not_found")
	ErrCodeSystemErr                                = errors.New("system_error")
	ErrCodeInvalidData                              = errors.New("invalid_data")
	ErrCodeItemNotInCart                            = errors.New("item_not_in_cart")
	ErrCodeItemNotFound                             = errors.New("item_not_found")
	ErrCodeInvalidCartStatus                        = errors.New("invalid_cart_status")
	ErrCodeInvalidItemQTY                           = errors.New("invalid_item_qty")
	ErrCodeInvalidTotalPrice                        = errors.New("invalid_total_price")
	ErrCodeYouAlreadyInClass                        = errors.New("you_already_in_class")
	ErrCodeHasInvalidItemInCart                     = errors.New("has_invalid_item_in_cart")
	ErrCodeClassClosed                              = errors.New("class_closed")
	ErrCodeClassNotFound                            = errors.New("class_not_found")
	ErrCodeUserDeactivated                          = errors.New("user_deactivated")
	ErrCodeTokenExpired                             = errors.New("token_expired")
	ErrCodeItemAlreadyExistsInOtherOrder            = errors.New("item_already_exists_in_other_order")
	ErrNotAuthorized                                = errors.New("not_authorized")
	ErrNotAllowTrialMode                            = errors.New("not_allow_trial_mode")
	ErrCodeGenerateExamFailed                       = errors.New("generate_exam_failed")
	ErrCodeYouAreNotInClass                         = errors.New("you_are_not_in_class")
	ErrCodeQuizNotFound                             = errors.New("quiz_not_found")
	ErrCodeRegistrationRequestNotFound              = errors.New("registration_request_not_found")
	ErrCodeUserAlreadyExists                        = errors.New("user_already_exists")
	ErrCodeUserNotExists                            = errors.New("user_not_exists")
	ErrCodeEmailAlreadyExists                       = errors.New("email_already_exists")
	ErrNotFoundRegistrationRequest                  = errors.New("not_found_registration_request")
	ErrInvalidRegistrationRequestPaymentStatus      = errors.New("invalid_registration_request_payment_status")
	ErrInvalidRegistrationRequestEntranceTestStatus = errors.New("invalid_registration_request_entrance_test_status")
	ErrInvalidRegistrationEnrollStatus              = errors.New("invalid_registration_request_enroll_status")
	ErrMissingEmail                                 = errors.New("missing_email")
	ErrActionNotAllowed                             = errors.New("action_not_allowed")
	ErrTokenNotFound                                = errors.New("token_not_found")
	ErrInvalidIeltsBand                             = errors.New("invalid_ielts_band")
	ErrInvalidStudyDuration                         = errors.New("invalid_study_duration")
	ErrCodeInvalidPromoCode                         = errors.New("invalid_promo_code")
	ErrCodePromoCodeIsExpired                       = errors.New("promo_code_is_expired")
	ErrCodePromoCodeIsUsedForSpecificItem           = errors.New("promo_code_is_used_for_specific_item")
	ErrCodePromoCodeIsUsedForMinAmount              = errors.New("promo_code_is_used_for_min_amount")
	ErrCodePromoCodeIsExceedMaxUse                  = errors.New("promo_code_is_exceed_max_use")
	ErrCodePromoCodeIsUsedForSpecificUser           = errors.New("promo_code_is_used_for_specific_user")
	ErrCodePromoCodeIsUsedForNewUser                = errors.New("promo_code_is_used_for_new_user")
	ErrCodeOrderNotFound                            = errors.New("order_not_found")
	ErrCodeInvalidOrderStatus                       = errors.New("invalid_order_status")
	ErrNotFoundEnv                                  = errors.New("not found env")
	ErrAlertFailed                                  = errors.New("alert failed")
	ErrNoQuestionsInQuiz                            = errors.New("no_questions_in_quiz")
	ErrRecordNotFound                               = errors.New("record_not_found")
	ErrQuizNotInClass                               = errors.New("quiz_not_in_class")
	ErrUserNotInClass                               = errors.New("user_not_in_class")
	ErrAnswerStatisticTypeRequired                  = errors.New("answer_statistic_type_required")
	ErrAnswerStatisticTypeInvalid                   = errors.New("answer_statistic_type_invalid")
	ErrNextExamDateBeforeNow                        = errors.New("next_exam_date_before_now")
	ErrInvalidVotes                                 = errors.New("invalid_votes")
	ErrNotFound                                     = gorm.ErrRecordNotFound
	ErrNoDataToUpdate                               = errors.New("no_data_to_update")
	ErrNotUpdateBothBandScoreAndBandScore100        = errors.New("not_update_both_band_score_and_band_score_100")
	ErrInvalidBandScore                             = errors.New("invalid_band_score")
	ErrInvalidAnswerStatus                          = errors.New("invalid_answer_status")
	ErrAnswerStatusReviewed                         = errors.New("answer_status_is_reviewed")
	ErrAnswerClassNull                              = errors.New("answer_class_is_null")
	ErrRecordExisted                                = errors.New("record_existed")
	ErrNoStudentsInClass                            = errors.New("no_students_in_class")
	ErrNoQuizzesInClass                             = errors.New("no_quizzes_in_class")
	ErrAnswerNotReviewedYet                         = errors.New("answer_not_reviewed_yet")
	ErrReachMaxRetries                              = errors.New("reach_max_retries")
	ErrReviewNotFound                               = errors.New("review_not_found")
	ErrInReviewing                                  = errors.New("in reviewing")
	ErrNotFoundRegistrations                        = errors.New("not_found_registration")
	ErrNotFoundStudent                              = errors.New("not_found_student")
	ErrNotFoundCourse                               = errors.New("not_found_course")
	ErrNoValidStudents                              = errors.New("no valid students")
	ErrNoValidClasses                               = errors.New("no valid classes")
	ErrNoValidCourses                               = errors.New("no valid courses")
	ErrNoValidRegistrationRequests                  = errors.New("no valid registration requests")
	ErrMaxAttachementExceeded                       = errors.New("max_attachment_exceeded")
	ErrInvalidObjectType                            = errors.New("invalid_object_type")
	ErrUpdateFailed                                 = errors.New("update_failed")
	ErrInvalidEmail                                 = errors.New("no_email")
	ErrInvalidRequest                               = errors.New("invalid_request")
	ErrInvalidAnswerID                              = errors.New("invalid Answer ID")
	ErrNoSlotAIReviewWriting                        = errors.New("student has run out of turns grading with AI")
	ErrInvalidWritingAnswer                         = errors.New("invalid writing answer")
	ErrNotFoundPosition                             = errors.New("not found position")
	ErrNoSlotAIVocabSuggest                         = errors.New("student has run out of turn using vocab suggestion")
	ErrInvalidCategory                              = errors.New("invalid category")
	ErrInvalidVocabValue                            = errors.New("invalid vocab value")
	ErrCategoryNotFound                             = errors.New("category not found in default categories")
	ErrVocabNotFound                                = errors.New("vocab entry not found or does not belong to the user")
	ErrCategoryNull                                 = errors.New("category must not be empty")
	ErrValueNull                                    = errors.New("value must not be empty")
	ErrNotRetrieveCategory                          = errors.New("cannot retrieve category")
	ErrInvalidVocabStatus                           = errors.New("invalid vocab status")
	ErrInvalidPromoCodeApplyContext                 = errors.New("invalid_promo_code_apply_context")
	ErrWrongLoginInfo                               = errors.New("wrong_login_info")
	ErrRegistrationRequestIsCancelled               = errors.New("registration request is cancelled")
	ErrCharacterNotFound                            = errors.New("character_not_found")
	ErrCharacterAlreadyExists                       = errors.New("character_already_exists")
)

var (
	ErrCodeInvalidTimeRange = errors.New("invalid_time_range")
)

var listErrorData = []errData{
	{
		Code:        "cart_not_found",
		HTTPCode:    404,
		MessageViVn: "Giỏ hàng không tồn tại",
		MessageEnUs: "Cart not found",
	},
	{
		Code:        "not_authorized",
		HTTPCode:    401,
		MessageViVn: "Không có quyền truy cập",
		MessageEnUs: "Not authorized",
	},
	{
		Code:        "invalid_cart_status",
		HTTPCode:    400,
		MessageViVn: "Trạng thái giỏ hàng không hợp lệ",
		MessageEnUs: "Invalid cart status",
	},
	{
		Code:        "item_not_in_cart",
		HTTPCode:    400,
		MessageViVn: "Sản phẩm không có trong giỏ hàng",
		MessageEnUs: "Item not in cart",
	},
	{
		Code:        "invalid_promo_code",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại không hợp lệ",
		MessageEnUs: "Invalid promo code",
	},
	{
		Code:        "exist_item_in_cart",
		HTTPCode:    400,
		MessageViVn: "Sản phẩm đã có trong giỏ hàng",
		MessageEnUs: "Exist item in cart",
	},
	{
		Code:        "invalid_item_qty",
		HTTPCode:    400,
		MessageViVn: "Số lượng sản phẩm không hợp lệ",
		MessageEnUs: "Invalid item qty",
	},
	{
		Code:        "has_invalid_item_in_cart",
		HTTPCode:    400,
		MessageViVn: "Giỏ hàng có sản phẩm không hợp lệ",
		MessageEnUs: "Cart has invalid item",
	},
	{
		Code:        "you_already_in_class",
		HTTPCode:    400,
		MessageViVn: "Bạn đã tham gia lớp học này rồi",
		MessageEnUs: "You already in class",
	},
	{
		Code:        "promo_code_is_used_for_specific_user",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại chỉ dành cho người dùng cụ thể",
		MessageEnUs: "Promo code is used for specific user",
	},
	{
		Code:        "promo_code_is_exceed_max_use",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại đã vượt quá số lần sử dụng",
		MessageEnUs: "Promo code is exceed max use",
	},
	{
		Code:        "promo_code_is_used_for_min_amount",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại chỉ dành cho đơn hàng có giá trị tối thiểu",
		MessageEnUs: "Promo code is used for min amount",
	},
	{
		Code:        "promo_code_is_used_for_specific_item",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại chỉ dành cho sản phẩm cụ thể",
		MessageEnUs: "Promo code is used for specific item",
	},
	{
		Code:        "promo_code_is_expired",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại đã hết hạn",
		MessageEnUs: "Promo code is expired",
	},
	{
		Code:        "invalid_time_range",
		HTTPCode:    400,
		MessageViVn: "Thời gian không hợp lệ",
		MessageEnUs: "Invalid time range",
	},
	{
		Code:        "invalid_payment_method",
		HTTPCode:    400,
		MessageViVn: "Phương thức thanh toán không hợp lệ",
		MessageEnUs: "Invalid payment method",
	},
	{
		Code:        "invalid_order_status",
		HTTPCode:    400,
		MessageViVn: "Trạng thái đơn hàng không hợp lệ",
		MessageEnUs: "Invalid order status",
	},
	{
		Code:        "invalid_data",
		HTTPCode:    400,
		MessageViVn: "Dữ liệu không hợp lệ",
		MessageEnUs: "Invalid data",
	},
	{
		Code:        "order_not_found",
		HTTPCode:    404,
		MessageViVn: "Đơn hàng không tồn tại",
		MessageEnUs: "Order not found",
	},
	{
		Code:        "invalid_total_price",
		HTTPCode:    400,
		MessageViVn: "Giá trị đơn hàng đã thay đổi, vui lòng kiểm tra lại!",
		MessageEnUs: "Total price is invalid",
	},
	{
		Code:        "system_error",
		HTTPCode:    500,
		MessageViVn: "Đã có lỗi xảy ra, vui lòng thử lại!",
		MessageEnUs: "System error, please try again!",
	},
	{
		Code:        "item_not_found",
		HTTPCode:    404,
		MessageViVn: "Sản phẩm không tồn tại",
		MessageEnUs: "Item not found",
	},
	{
		Code:        "class_closed",
		HTTPCode:    400,
		MessageViVn: "Lớp học đã hết hạn",
		MessageEnUs: "Class closed",
	},
	{
		Code:        "class_not_found",
		HTTPCode:    404,
		MessageViVn: "Khóa học không tồn tại hoặc không khả dụng",
		MessageEnUs: "Class not found",
	},
	{
		Code:        "user_deactivated",
		HTTPCode:    403,
		MessageViVn: "Tài khoản của bạn đã bị vô hiệu hóa",
		MessageEnUs: "Your account has been deactivated",
	},
	{
		Code:        "promo_code_is_used_for_new_user",
		HTTPCode:    400,
		MessageViVn: "Mã khuyến mại chỉ dành cho người dùng mới",
		MessageEnUs: "Promo code is used for new user",
	},
	{
		Code:        "token_expired",
		HTTPCode:    403,
		MessageViVn: "Phiên làm việc của bạn đã hết hạn",
		MessageEnUs: "Your session has expired",
	},
	{
		Code:        "item_already_exists_in_other_order",
		HTTPCode:    400,
		MessageViVn: "Đã tồn tại sản phẩm này trong 1 đơn hàng khác, vui lòng kiểm tra lại!",
		MessageEnUs: "This item is already in another order, please check again!",
	},
	{
		Code:        "not_allow_trial_mode",
		HTTPCode:    400,
		MessageViVn: "Lớp học không hỗ trợ dạng học thử!",
		MessageEnUs: "Class not support trial mode!",
	},
	{
		Code:        "not_found_registration_request",
		HTTPCode:    400,
		MessageViVn: "Không tìm thấy đơn phù hợp",
		MessageEnUs: "Not found valid request",
	},
	{
		Code:        "invalid_registration_request_payment_status",
		HTTPCode:    400,
		MessageViVn: "Trạng thái thanh toán đã hoàn thành hoặc không phù hợp",
		MessageEnUs: "Payment status completed or invalid",
	},
	{
		Code:        "invalid_registration_request_entrance_test_status",
		HTTPCode:    400,
		MessageViVn: "Trạng thái bài test đã hoàn thành hoặc không phù hợp",
		MessageEnUs: "Entrance test status completed or invalid",
	},
	{
		Code:        "invalid_registration_request_enroll_status",
		HTTPCode:    400,
		MessageViVn: "Đã tham gia vào lớp học hoặc trạng thái không phù hợp",
		MessageEnUs: "Enroll status completed or invalid",
	},
	{
		Code:        "registration_request_not_found",
		HTTPCode:    404,
		MessageViVn: "Yêu cầu đăng ký không tồn tại",
		MessageEnUs: "Registration request not found",
	},
	{
		Code:        "user_already_exists",
		HTTPCode:    400,
		MessageViVn: "Tài khoản đã tồn tại",
		MessageEnUs: "User already exists",
	},
	{
		Code:        "user_not_exists",
		HTTPCode:    404,
		MessageViVn: "Tài khoản không tồn tại",
		MessageEnUs: "User not exists",
	},
	{
		Code:        "email_already_exists",
		HTTPCode:    400,
		MessageViVn: "Email đã tồn tại",
		MessageEnUs: "Email already exists",
	},
	{
		Code:        "not_found_registration_request",
		HTTPCode:    400,
		MessageViVn: "Không tìm thấy đơn phù hợp",
		MessageEnUs: "Not found valid request",
	},
	{
		Code:        "invalid_registration_request_payment_status",
		HTTPCode:    400,
		MessageViVn: "Trạng thái thanh toán đã hoàn thành hoặc không phù hợp",
		MessageEnUs: "Payment status completed or invalid",
	},
	{
		Code:        "invalid_registration_request_entrance_test_status",
		HTTPCode:    400,
		MessageViVn: "Trạng thái bài test đã hoàn thành hoặc không phù hợp",
		MessageEnUs: "Entrance test status completed or invalid",
	},
	{
		Code:        "invalid_registration_request_enroll_status",
		HTTPCode:    400,
		MessageViVn: "Đã tham gia vào lớp học hoặc trạng thái không phù hợp",
		MessageEnUs: "Enroll status completed or invalid",
	},
	{
		Code:        "missing_email",
		HTTPCode:    400,
		MessageViVn: "Email rỗng",
		MessageEnUs: "Empty email",
	},
	{
		Code:        "action_not_allowed",
		HTTPCode:    403,
		MessageViVn: "Bạn không có quyền thực hiện",
		MessageEnUs: "You're not allow",
	},
	{
		Code:        "token_not_found",
		HTTPCode:    401,
		MessageViVn: "Không tìm thấy token",
		MessageEnUs: "Token not found",
	},
	{
		Code:        "invalid_ielts_band",
		HTTPCode:    400,
		MessageViVn: "Điểm số ielts không hợp lệ",
		MessageEnUs: "Invalid band",
	},
	{
		Code:        "invalid_study_duration",
		HTTPCode:    400,
		MessageViVn: "Số giờ học không hợp lệ",
		MessageEnUs: "Invalid study duration",
	},
	{
		Code:        ErrRecordNotFound.Error(),
		HTTPCode:    http.StatusNotFound,
		MessageViVn: "Record not found",
		MessageEnUs: "Record not found",
	},
	{
		Code:        ErrQuizNotInClass.Error(),
		HTTPCode:    http.StatusNotFound,
		MessageViVn: "Quiz không thuộc lớp học",
		MessageEnUs: "Quiz is not in the class",
	},
	{
		Code:        ErrUserNotInClass.Error(),
		HTTPCode:    http.StatusNotFound,
		MessageViVn: "Người dùnng không thuộc lớp học",
		MessageEnUs: "User is not in the class",
	},
	{
		Code:        ErrAnswerStatisticTypeRequired.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Type là bắt buộc, vui lòng nhập thông tin cho field này",
		MessageEnUs: "Required for type",
	},
	{
		Code:        ErrAnswerStatisticTypeInvalid.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Type không hợp lệ, vui lòng nhập thông tin hợp lệ",
		MessageEnUs: "Type is invalid",
	},
	{
		Code:        ErrNextExamDateBeforeNow.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Thời gian thi phải sau thời gian hiện tại",
		MessageEnUs: "Exam date must be after now",
	},
	{
		Code:        ErrNotFound.Error(),
		HTTPCode:    http.StatusNotFound,
		MessageViVn: "Không tìm thấy",
		MessageEnUs: "Not found",
	},
	{
		Code:        ErrNoDataToUpdate.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Không có thông tin để cập nhật",
		MessageEnUs: "No data to update",
	},
	{
		Code:        ErrNotUpdateBothBandScoreAndBandScore100.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Không thể cập nhật đồng thời điểm theo ielts và điểm theo thang 100",
		MessageEnUs: "Can't update both band score and band score with scale 100",
	},
	{
		Code:        ErrInvalidBandScore.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Điểm số không hợp lệ",
		MessageEnUs: "Band score/band score with scale 100 is invalid",
	},
	{
		Code:        ErrInvalidAnswerStatus.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Trạng thái của câu trả lời không hợp lệ",
		MessageEnUs: "Answer status is invalid",
	},
	{
		Code:        ErrAnswerStatusReviewed.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Câu trả lời đã được chấm điểm, vui lòng không cập nhật thông tin",
		MessageEnUs: "Answer status is invalid",
	},
	{
		Code:        ErrAnswerClassNull.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Lớp học ở câu trả lời không có thông tin",
		MessageEnUs: "Answer class is null",
	},
	{
		Code:        ErrActionNotAllowed.Error(),
		HTTPCode:    http.StatusForbidden,
		MessageViVn: "Không có quyền thực hiện thao tác",
		MessageEnUs: "Action is not allowed",
	},
	{
		Code:        ErrRecordExisted.Error(),
		HTTPCode:    http.StatusConflict,
		MessageViVn: "Record đã tồn tại",
		MessageEnUs: "Record is already existed",
	},
	{
		Code:        ErrNoStudentsInClass.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Lớp không có học sinh",
		MessageEnUs: "Class doesn't have any student",
	},
	{
		Code:        ErrNoQuizzesInClass.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Lớp không có các quiz hợp lệ",
		MessageEnUs: "Class doesn't have any valid quizzes",
	},
	{
		Code:        ErrAnswerNotReviewedYet.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Bài làm chưa được review",
		MessageEnUs: "Answer hasn't been reviewed",
	},
	{
		Code:        ErrReviewNotFound.Error(),
		HTTPCode:    http.StatusNotFound,
		MessageViVn: "Review không tồn tại",
		MessageEnUs: "Review is not found",
	},
	{
		Code:        ErrMaxAttachementExceeded.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Số lượng file đính kèm vượt quá giới hạn",
		MessageEnUs: "Number of attachments exceeds the limit",
	},
	{
		Code:        ErrInvalidObjectType.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Loại đối tượng không hợp lệ",
		MessageEnUs: "Invalid object type",
	},
	{
		Code:        ErrUpdateFailed.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Cập nhật thất bại",
		MessageEnUs: "Update failed",
	},
	{
		Code:        ErrInvalidEmail.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Email không hợp lệ",
		MessageEnUs: "Invalid email",
	},
	{
		Code:        ErrInvalidRequest.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Yêu cầu không hợp lệ",
		MessageEnUs: "Invalid request",
	},
	{
		Code:        ErrNotFoundPosition.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Vị trí không hợp lệ",
		MessageEnUs: "Invalid position",
	},
	{
		Code:        ErrInvalidPromoCodeApplyContext.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Mã khuyến mại không áp dụng cho đối tượng này",
		MessageEnUs: "Promo code is not applied for this context",
	},
	{
		Code:        ErrWrongLoginInfo.Error(),
		HTTPCode:    http.StatusUnauthorized,
		MessageViVn: "Thông tin đăng nhập không đúng",
		MessageEnUs: "Login information is wrong",
	},
	{
		Code:        ErrRegistrationRequestIsCancelled.Error(),
		HTTPCode:    http.StatusBadRequest,
		MessageViVn: "Yêu cầu đăng ký đã bị hủy",
		MessageEnUs: "Registration request is cancelled",
	},
}

var (
	AllErrors *MasterErrData
)

func FetchMasterErrData() {
	AllErrors = NewMasterErrData()
	AllErrors.fetchAll()
}

type errData struct {
	Code        string `json:"code" gorm:"column:code"`
	HTTPCode    int    `json:"http_code" gorm:"column:http_code"`
	MessageViVn string `json:"message_vi_vn" gorm:"column:message_vi_vn"`
	MessageEnUs string `json:"message_en_us" gorm:"column:message_en_us"`
}

type ExtraData struct {
	OrderID int64 `json:"order_id,omitempty"`
}

type LocalizeErrRes struct {
	Code      string     `json:"code,omitempty"`
	Message   string     `json:"message,omitempty"`
	HTTPCode  int        `json:"-"`
	Internal  string     `json:"internal,omitempty"`
	ExtraData *ExtraData `json:"extra_data,omitempty"`
}

func (a *LocalizeErrRes) Error() string {
	return a.Code
}

type MasterErrData struct {
	mutex sync.Mutex
	data  map[string]errData
}

// Error data

func NewMasterErrData() *MasterErrData {
	return &MasterErrData{}
}

func (a *MasterErrData) fetchAll() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for _, errMessage := range listErrorData {
		if a.data == nil {
			a.data = make(map[string]errData)
		}
		a.data[errMessage.Code] = errMessage
	}
}

func (a *MasterErrData) New(err error, language string, internal ...string) *LocalizeErrRes {
	errRes := new(LocalizeErrRes)
	ok := errors.As(err, &errRes)
	if !ok {
		errRes = &LocalizeErrRes{
			Code:    "bad_request",
			Message: "Đã có lỗi xảy ra, vui lòng thử lại!",
		}
		if len(internal) > 0 {
			errRes.Internal = internal[0]
		}
		errFromDB, exists := a.data[err.Error()]
		if exists {
			errRes.Code = errFromDB.Code
			errRes.HTTPCode = errFromDB.HTTPCode
			switch language {
			case "vi":
				errRes.Message = errFromDB.MessageViVn
			default:
				errRes.Message = errFromDB.MessageEnUs
			}
		} else {
			errRes.HTTPCode = 400
		}
	}

	if len(internal) > 0 {
		errRes.Internal = internal[0]
	}
	return errRes
}

// Error res

func (a *LocalizeErrRes) SetMessage(message string) *LocalizeErrRes {
	a.Message = message
	return a
}

func (a *LocalizeErrRes) ReplaceDescByVars(args ...interface{}) *LocalizeErrRes {
	for _, arg := range args {
		a.Message = fmt.Sprintf(a.Message, arg)
	}
	return a
}

func (a *LocalizeErrRes) SetOrderIDToExtraData(orderID int64) *LocalizeErrRes {
	if a.ExtraData == nil {
		a.ExtraData = new(ExtraData)
	}
	a.ExtraData.OrderID = orderID
	return a
}

func (a *LocalizeErrRes) ConvertToBaseError() Response {
	res := BaseResponse(REQUEST_FAILED, a.Message, a.Internal, a.ExtraData)
	res.SetErrorCode(a.Code)
	return res
}

func AbortWithError(c *gin.Context, err error) {
	errJSON := AllErrors.New(err, "vi", err.Error())
	c.AbortWithStatusJSON(errJSON.HTTPCode, errJSON.ConvertToBaseError())
}
