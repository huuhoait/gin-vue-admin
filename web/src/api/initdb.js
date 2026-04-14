import service from '@/utils/request'
// @Tags InitDB
// @Summary Initialize database
// @Produce  application/json
// @Param data body request.InitDB true "Database initialization parameters"
// @Success 200 {string} string "{"code":0,"data":{},"msg":"Database created"}"
// @Router /init/initdb [post]
export const initDB = (data) => {
  return service({
    url: '/init/initdb',
    method: 'post',
    data,
    donNotShowLoading: true
  })
}

// @Tags CheckDB
// @Summary Check database
// @Produce  application/json
// @Success 200 {string} string "{"code":0,"data":{},"msg":"OK"}"
// @Router /init/checkdb [post]
export const checkDB = () => {
  return service({
    url: '/init/checkdb',
    method: 'post'
  })
}
