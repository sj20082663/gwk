package ggk

import (
	"container/list"
	"runtime"
)

type CanvasInitFlags int

const (
	KCanvasInitFlagDefault = 1 << iota
	KCanvasInitFlagConservativeRasterClip
)

type Canvas struct {
	surfaceProps SurfaceProps
	saveCount    int
	metaData     *MetaData
	baseSurface  *BaseSurface
	mcStack      *list.List
	clipStack    *ClipStack
	mcRec        *tCanvasMCRec // points to top of the stack

	deviceCMDirty bool

	cachedLocalClipBounds      Rect
	cachedLocalClipBoundsDirty bool

	allowSoftClip          bool
	allowSimplifyClip      bool
	conservativeRasterClip bool
}

func NewCanvas(bmp *Bitmap) *Canvas {
	var canvas = new(Canvas)

	canvas.surfaceProps = MakeSurfaceProps(KSurfacePropsFlagNone,
		KSurfacePropsInitTypeLegacyFontHost)
	canvas.mcStack = list.New()
	var device = NewBitmapDevice(bmp, canvas.surfaceProps)
	canvas.init(device.BaseDevice, KCanvasInitFlagDefault)
	return canvas
}

func (c *Canvas) init(dev *BaseDevice, flags CanvasInitFlags) {
	toimpl()
}

func (c *Canvas) ReadPixelsToBitmap(bmp *Bitmap, x, y Scalar) error {
	toimpl()
	return nil
}

func (c *Canvas) ReadPixelsInRectToBitmap(bmp *Bitmap, srcRect Rect) error {
	toimpl()
	return nil
}

func (c *Canvas) ReadPixels(dstInfo *ImageInfo, dstData []byte, rowBytes int,
	x, y Scalar) error {

	var dev = c.Device()
	if dev == nil {
		return errorf("device is nil")
	}

	var size = c.BaseLayerSize()

	var hlp = newReadPixelsHlp(dstInfo, dstData, rowBytes, x, y)
	if err := hlp.Trim(size.Width(), size.Height()); err != nil {
		return errorf("bad param %v", err)
	}

	// The device can assert that the requested area is always contained in its
	// bounds.
	return dev.ReadPixels(hlp.Info, hlp.Pixels, hlp.RowBytes, hlp.X, hlp.Y)
}

func (c *Canvas) Device() *BaseDevice {
	var rec = c.mcStack.Front().Value.(*tCanvasMCRec)
	return rec.Layer.Device
}

func (c *Canvas) BaseLayerSize() Size {
	var dev = c.Device()
	var sz Size
	if dev != nil {
		sz = MakeSize(dev.Width(), dev.Height())
	}
	return sz
}

type CanvasSaveFlags int

const (
	KSaveFlagHasAlphaLayer   = 0x01
	KSaveFlagFullColorLayer  = 0x02
	KSaveFlagClipToLayer     = 0x10
	KSaveFlagARGBNoClipLayer = 0x0F
	KSaveFlagARGBClipLayer   = 0x1F
)

// tDeviceCM is the record we keep for each BaseDevice that the user installs.
// The clip/matrix/proc are fields that reflect the top of the save/resotre
// stack. Whenever the canvas changes, it makes a dirty flag, and then before
// these are used (assuming we're not on a layer) we rebuild these cache values:
// they reflect the top of the save stack, but translated and clipped by the
// device's XY offset and bitmap-bounds.
type tDeviceCM struct {
	Next                 *tDeviceCM
	Device               *BaseDevice
	Clip                 *RasterClip
	Paint                *Paint
	Matrix               *Matrix
	MatrixStroage        *Matrix
	DeviceIsBitmapDevice bool
}

func newDeivceCM(dev *BaseDevice, paint *Paint, canvas *Canvas,
	conservativeRasterClip bool, deviceIsBitmapDevice bool) {

	var deviceCM = new(tDeviceCM)
	deviceCM.Next = nil
	deviceCM.Clip = NewRasterClip(conservativeRasterClip)
	deviceCM.DeviceIsBitmapDevice = deviceIsBitmapDevice
	if dev != nil {
		dev.OnAttachToCanvas(canvas)
	}
	runtime.SetFinalizer(deviceCM, func(d *tDeviceCM) {
		d.Device.OnDetachFromCanvas()
	})
	deviceCM.Device = dev
	deviceCM.Paint = paint
}

func (d *tDeviceCM) Reset(bounds Rect) {
	d.Clip.SetRect(bounds)
}

func (d *tDeviceCM) UpdateMC(totalMatrix *Matrix, totlaClip *RasterClip,
	clipStack *ClipStack, updateClip *RasterClip) {

	toimpl()
}

// tCanvasMCRec is the record we keep for each save/restore level in the stack.
// Since a level optionally copies the matrix and/or stack, we have pointers
// for these fields. If the value is copied for this level, the copy is stored
// in the ...Storage field, and the pointer points to that. If the value is not
// copied for this level, we ignore ...Storage, and just point at the
// corresponding value in the previous level in the stack.
type tCanvasMCRec struct {
	Layer *Layer
}
