import service from '@/utils/request'
// @Tags email
// @Summary Send test email
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"OK"}"
// @Router /email/emailTest [post]
export const emailTest = (data) => {
  return service({
    url: '/email/emailTest',
    method: 'post',
    data
  })
}
